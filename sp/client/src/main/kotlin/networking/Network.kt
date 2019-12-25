package networking

import javafx.application.Platform
import javafx.scene.control.Alert
import model.lobby.User
import mu.KotlinLogging
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.ErrorResponseMessage
import networking.messages.GetLobbiesMessage
import networking.reader.SimpleMessageReader
import networking.receiver.FixedMessageReceiver
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
        get() {
            synchronized(field) {
                return field
            }
        }

    private var responseListeners = mutableMapOf<Int, MutableList<(ApplicationMessage) -> Unit>>()

    val connectionStatusListeners = mutableListOf<ConnectionStatusListener>()

    fun connectTo(hostname: String, port: Int) {
        tcpLayer?.close()
        tcpLayer = TCPLayer(port, hostname, FixedMessageReceiver(SimpleMessageReader(this)), this)
        tcpLayer?.start()
    }

    fun removeMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        KotlinLogging.logger {  }.debug { "Adding a message listener of for messages of class ${messageClazz.simpleName}" }
        messageListeners[messageClazz]?.remove(callback)
    }

    inline fun <reified T : ApplicationMessage> addMessageListener(noinline callback: (T) -> Unit) {
        KotlinLogging.logger {  }.debug { "Adding a message listener of for messages of class ${T::class.java.simpleName}" }
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
        logger.info { "Received a message of tyaape ${message.type}" }
        synchronized(messageListeners) {
            val typeMessageListeners = messageListeners.getOrDefault(message.javaClass, listOf<(ApplicationMessage) -> Unit>())
            logger.info { "About to invoke ${typeMessageListeners.size} message listeners of type ${message.type}" }
            for (callback: (ApplicationMessage) -> Unit in typeMessageListeners) {
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
        val json = message.toJson().replace(FixedMessageReceiver.START_CHAR.toString(), "\\" + FixedMessageReceiver.START_CHAR).replace(FixedMessageReceiver.SEPARATOR.toString(), "\\" + FixedMessageReceiver.SEPARATOR)

        if (message is GetLobbiesMessage) {
            println()
        }

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

        tcpLayer?.write("${FixedMessageReceiver.START_CHAR}${json.toByteArray().size}${FixedMessageReceiver.SEPARATOR}${message.type}${FixedMessageReceiver.SEPARATOR}${messageId}${FixedMessageReceiver.SEPARATOR}$json")

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