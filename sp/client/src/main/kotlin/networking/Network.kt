package networking

import javafx.application.Platform
import javafx.scene.control.Alert
import model.lobby.User
import mu.KotlinLogging
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.ErrorResponseMessage
import networking.reader.SimpleMessageReader
import networking.receiver.SimpleMessageReceiver
import tornadofx.alert

private val logger = KotlinLogging.logger { }

class Network : ConnectionStatusListener, ApplicationMessageReader {

    companion object {
        private val MESSAGE_ID_CEILING = 60000
        private val MESSAGE_STARTING_ID = 1

        private var network: Network = Network()

        public lateinit var User: User

        @Synchronized
        fun getInstance(): Network {
            return network
        }


    }

    private var messageId = MESSAGE_STARTING_ID

    var tcpLayer: TCPLayer? = null

    var messageListeners = mutableMapOf<Class<out ApplicationMessage>, MutableList<(T: ApplicationMessage) -> Unit>>()

    private var responseListeners = mutableMapOf<Int, MutableList<(ApplicationMessage) -> Unit>>()

    val connectionStatusListeners = mutableListOf<ConnectionStatusListener>()

    fun connectTo(hostname: String, port: Int) {
        tcpLayer?.close()
        tcpLayer = TCPLayer(port, hostname, SimpleMessageReceiver(SimpleMessageReader(this)), this)
        tcpLayer?.start()
    }

    fun removeMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        messageListeners[messageClazz]?.remove(callback)
    }

    inline fun <reified T : ApplicationMessage> addMessageListener(noinline callback: (T) -> Unit) {
        messageListeners.putIfAbsent(T::class.java, mutableListOf())
        messageListeners[T::class.java]?.add(callback as (ApplicationMessage) -> Unit)
    }

    inline fun <reified T : ApplicationMessage> removeMessageListener(noinline callback: (T) -> Unit) {
        messageListeners[T::class.java]?.remove(callback)
    }

    fun addResponseListener(mid: Int, callback: (ApplicationMessage) -> Unit) {
        responseListeners.putIfAbsent(mid, mutableListOf())
        responseListeners[mid]?.add(callback)
    }

    fun removeResponseListener(mid: Int, callback: (ApplicationMessage) -> Unit) {
        responseListeners[mid]?.remove(callback)
    }

    override fun onConnected() {
        connectionStatusListeners.forEach {
            it.onConnected()
        }
    }

    override fun onUnreachable() {
        connectionStatusListeners.forEach {
            it.onUnreachable()
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        connectionStatusListeners.forEach {
            it.onFailedAttempt(attempt)
        }
    }

    override fun read(message: ApplicationMessage, mid: Int) {
        logger.info { "Received a message of type ${message.type}" }
        synchronized(messageListeners) {
            logger.info { "About to invoke ${messageListeners.size} message listeners of type ${message.type}" }
            for (callback: (ApplicationMessage) -> Unit in messageListeners.getOrDefault(message.javaClass, listOf<(ApplicationMessage) -> Unit>())) {
                callback.invoke(message)
            }
        }

        synchronized(responseListeners) {
            val idResponseListeners = responseListeners.getOrDefault(mid, listOf<(ApplicationMessage) -> Unit>())
            logger.info { "About to invoke ${idResponseListeners.size} response listeners of message ID ${mid} and type ${message.type}" }
            for (callback: (ApplicationMessage) -> Unit in idResponseListeners) {
                callback.invoke(message)
            }
            responseListeners[mid]?.clear()
        }
    }

    fun send(message: ApplicationMessage, callback: ((ApplicationMessage) -> Unit)? = null, desiredMessageId: Int = 0, callAfterWrite: (() -> Unit)? = null, ignoreErrors: Boolean = false) {
        val json = message.toJson()

        callback?.let {
            addResponseListener(if (desiredMessageId != 0) desiredMessageId else messageId, callback)
            if (!ignoreErrors) {
                addResponseListener(if (desiredMessageId != 0) desiredMessageId else messageId) { am: ApplicationMessage ->
                    run {
                        if (am is ErrorResponseMessage)
                            Platform.runLater {
                                alert(Alert.AlertType.ERROR, "Error", am.content)
                            }
                    }
                }
            }
        }

        logger.info { "Printing message of type ${message.type} and content '$json' to server" }

        tcpLayer?.write("${SimpleMessageReceiver.START_CHAR}${json.toByteArray().size}${SimpleMessageReceiver.SEPARATOR}${message.type}${SimpleMessageReceiver.SEPARATOR}${messageId}${SimpleMessageReceiver.SEPARATOR}$json")

        messageId++

        if (messageId > MESSAGE_ID_CEILING) {
            messageId = MESSAGE_STARTING_ID
        }

        callAfterWrite?.invoke()
    }

    fun stop() {
        tcpLayer?.close()
    }
}