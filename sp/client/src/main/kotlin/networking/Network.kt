package networking

import javafx.application.Platform
import javafx.scene.control.Alert
import model.lobby.User
import mu.KotlinLogging
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.*
import networking.reader.SimpleMessageReader
import networking.receiver.FixedMessageReceiver
import screens.*
import screens.game.GameStateRegenerationEvent
import tornadofx.FX
import tornadofx.alert
import java.time.LocalDateTime
import java.util.*
import kotlin.concurrent.timer

private val logger = KotlinLogging.logger { }

class Network : ConnectionStatusListener, ApplicationMessageReader {

    companion object {
        private val MESSAGE_ID_CEILING = 60000
        private val MESSAGE_STARTING_ID = 1
        private const val RESPONSE_CLEANER_TIMER = 2000L
        private const val REAUTH_ATTEMPT_CEILING = 3

        private var network: Network = Network()

        public lateinit var User: User
        public var authorized = false

        @Synchronized
        fun getInstance(): Network {
            return network
        }


    }

    private var messageId = MESSAGE_STARTING_ID

    var tcpLayer: TCPLayer? = null
    private var responseCleanerTimer: Timer? = null
    private var keepAliveTimer: Timer? = null
    private var triedToReconnect = false
    private var connected = false
    private var keepAliveSent = false

    var messageListeners = mutableMapOf<Class<out ApplicationMessage>, MutableList<(T: ApplicationMessage) -> Unit>>()
        get() {
            synchronized(field) {
                return field
            }
        }

    private var responseListeners = mutableMapOf<Int, MutableList<ResponseCallback>>()
        get() {
            synchronized(field) {
                return field
            }
        }

    val connectionStatusListeners = mutableListOf<ConnectionStatusListener>()

    lateinit var hostname: String
    var port: Int = 0

    fun connectTo(hostname: String, port: Int) {
        this.hostname = hostname
        this.port = port
        tcpLayer?.close()
        tcpLayer = TCPLayer(port, hostname, FixedMessageReceiver(SimpleMessageReader(this)), this)
        tcpLayer?.start()
    }

    fun removeMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        KotlinLogging.logger { }.debug { "Adding a message listener of for messages of class ${messageClazz.simpleName}" }
        messageListeners[messageClazz]?.remove(callback)
    }

    inline fun <reified T : ApplicationMessage> addMessageListener(noinline callback: (T) -> Unit) {
        KotlinLogging.logger { }.debug { "Adding a message listener of for messages of class ${T::class.java.simpleName}" }
        messageListeners.putIfAbsent(T::class.java, mutableListOf())
        messageListeners[T::class.java]?.add(callback as (ApplicationMessage) -> Unit)
    }

    inline fun <reified T : ApplicationMessage> removeMessageListener(noinline callback: (T) -> Unit) {
        messageListeners[T::class.java]?.remove(callback)
    }

    fun addResponseListener(mid: Int, callback: ResponseCallback) {
        synchronized(responseListeners) {
            responseListeners.putIfAbsent(mid, mutableListOf())
            responseListeners[mid]?.add(callback)
        }
    }

    fun removeResponseListener(mid: Int, callback: ResponseCallback) {
        responseListeners[mid]?.remove(callback)
    }

    override fun onConnected() {
        logger.info { "Network OnConnected" }
        connected = true
        keepAliveSent = false
        if (!triedToReconnect) {
            connectionStatusListeners.forEach {
                it.onConnected()
            }
        } else {
            connectionStatusListeners.forEach {
                it.onReconnected()
            }
        }

        logger.info { "triedToReconnect:$triedToReconnect, authorized:$authorized" }
        if (triedToReconnect && authorized) {
            logger.info { "Sending authentication" }
            reauthAttempt = 0
            tryToReauthenticate()
        }

        triedToReconnect = false

        responseCleanerTimer?.cancel()
        keepAliveTimer?.cancel()
        responseCleanerTimer = timer(period = RESPONSE_CLEANER_TIMER) {
            synchronized(responseListeners) {
                val messageIDsWhosListenersToBeRemoved = mutableListOf<Int>()
                responseListeners.forEach { listeners ->
                    logger.info { "There are ${listeners.value.size} listeners for mid ${listeners.key}" }
                    val listenersToBeRemoved = mutableListOf<ResponseCallback>()
                    listeners.value.forEach {
                        if (it.timestamp.plusSeconds(2).isBefore(LocalDateTime.now())) {
                            listenersToBeRemoved.add(it)
                            it.timeoutCallBack?.invoke()
                        }
                    }
                    listeners.value.removeAll(listenersToBeRemoved)
                    if (listeners.value.size <= 0) {
                        messageIDsWhosListenersToBeRemoved.add(listeners.key)
                    }
                }
                messageIDsWhosListenersToBeRemoved.forEach { responseListeners.remove(it) }
            }
        }

        keepAliveSent = false
        keepAliveTimer = timer(period = RESPONSE_CLEANER_TIMER) {
            if (!keepAliveSent) {
                keepAliveSent = true
                send(KeepAliveMessage(), callback = {
                    keepAliveSent = false
                }, ignoreErrors = true, timeoutCallback = {
                    onKeepAliveFailed()
                    keepAliveTimer?.cancel()
                })
            }
        }
    }

    var reauthAttempt = 0

    private fun tryToReauthenticate() {
        if (reauthAttempt < REAUTH_ATTEMPT_CEILING) {
            reauthAttempt++
            Thread.sleep(2000)
            send(UserAuthenticationMessage(User.name, reconnecting = true), {
                when (it) {
                    is UserStateRegenerationResponse -> {
                        when (it.state) {
                            UserStateRegenerationResponse.SERVER_RESTARTED -> {
                                User = it.user!!
                                authorized = true
                                FX.eventbus.fire(ServerRestartedEvent())
                            }
                            UserStateRegenerationResponse.SERVER_RESTARTED_NAME_TAKEN -> {
                                authorized = false
                                tryToReauthenticate()
                            }
                            UserStateRegenerationResponse.GAME -> {
                                authorized = true
                                Platform.runLater {
                                    alert(Alert.AlertType.ERROR, "This should not have happened")
                                }
                            }
                        }
                    }

                    is GameStateRegenerationResponse -> {
                        authorized = true
                        User = it.user
                        FX.eventbus.fire(GameRegeneratedEvent(it))
                    }

                    is ErrorResponseMessage -> {
                        authorized = false
                        tryToReauthenticate()
                    }
                    else -> {
                        Platform.runLater {
                            alert(Alert.AlertType.ERROR, "Failed to reauthenticate for unknown reasons")
                        }
                    }
                }
            }, ignoreErrors = true)
        } else {
            Platform.runLater {
                alert(Alert.AlertType.ERROR, "Failed to reauthenticate after reconnect")
            }
            FX.eventbus.fire(ServerRestartedUnauthorizedEvent())
        }
    }

    override fun onUnreachable() {
        logger.warn { "onUnreachable" }
        connected = false
        responseCleanerTimer?.cancel()
        keepAliveTimer?.cancel()
        FX.eventbus.fire(ServerUnreachableEvent())
        connectionStatusListeners.forEach {
            it.onUnreachable()
        }
        triedToReconnect = false
    }

    @Synchronized
    private fun reconnect() {
        logger.info { "Trying to reconnect triedToReconnect=$triedToReconnect" }
        if (!triedToReconnect) {
            triedToReconnect = true
            tcpLayer?.close()
            tcpLayer = TCPLayer(port, hostname, FixedMessageReceiver(SimpleMessageReader(this)), this)
            tcpLayer?.start()
        }
    }

    override fun onReconnected() {
        connectionStatusListeners.forEach {
            it.onReconnected()
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
            val typeMessageListeners = messageListeners.getOrDefault(message.javaClass, listOf<(ApplicationMessage) -> Unit>())
            logger.info { "About to invoke ${typeMessageListeners.size} message listeners of type ${message.type}" }
            for (callback: (ApplicationMessage) -> Unit in typeMessageListeners) {
                callback.invoke(message)
            }
        }

        synchronized(responseListeners) {
            val idResponseListeners = responseListeners.getOrDefault(mid, listOf<ResponseCallback>())
            logger.info { "About to invoke ${idResponseListeners.size} response listeners of message ID ${mid} and type ${message.type}" }
            for (callback: ResponseCallback in idResponseListeners) {
                callback.callback.invoke(message)
            }
            responseListeners[mid]?.clear()
            responseListeners.remove(mid)

        }
    }

    fun send(message: ApplicationMessage, callback: ((ApplicationMessage) -> Unit)? = null, desiredMessageId: Int = 0, callAfterWrite: (() -> Unit)? = null, ignoreErrors: Boolean = false, timeoutCallback: (() -> Unit)? = null) {
        val json = message.toJson().replace(FixedMessageReceiver.START_CHAR.toString(), "\\" + FixedMessageReceiver.START_CHAR).replace(FixedMessageReceiver.SEPARATOR.toString(), "\\" + FixedMessageReceiver.SEPARATOR)

        if (message is GetLobbiesMessage) {
            println()
        }

        callback?.let {
            addResponseListener(if (desiredMessageId != 0) desiredMessageId else messageId, ResponseCallback(callback, timeoutCallback))
            if (!ignoreErrors) {
                addResponseListener(if (desiredMessageId != 0) desiredMessageId else messageId, ResponseCallback({ am: ApplicationMessage ->
                    run {
                        if (am is ErrorResponseMessage)
                            Platform.runLater {
                                alert(Alert.AlertType.ERROR, "Error", am.content)
                            }
                    }
                }, timeoutCallBack = {
                    Platform.runLater {
                        alert(Alert.AlertType.ERROR, "Error", "Operation timeout for message ${message.javaClass}")
                    }
                }))
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

    private fun onKeepAliveFailed() {
        FX.eventbus.fire(DisconnectedEvent())
        logger.info { "From keepAlive called reconnect" }
        reconnect()
    }

    fun stop() {
        tcpLayer?.close()
    }
}

class ResponseCallback(val callback: (ApplicationMessage) -> Unit, val timeoutCallBack: (() -> Unit)? = null) {
    internal val timestamp = LocalDateTime.now()
}