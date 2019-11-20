package networking

import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.SuccessResponseMessage
import networking.reader.SimpleMessageReader
import networking.receiver.SimpleMessageReceiver
import tornadofx.App

class Network : ConnectionStatusListener, ApplicationMessageReader {

    companion object {

        private val MESSAGE_ID_CEILING = 60000
        private val MESSAGE_STARTING_ID = 1

        private var network: Network = Network()

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

    inline fun <reified T: ApplicationMessage> addMessageListener(noinline callback: (T) -> Unit){
        messageListeners.putIfAbsent(T::class.java, mutableListOf())
        messageListeners[T::class.java]?.add(callback as (ApplicationMessage) -> Unit)
    }

    inline fun <reified T: ApplicationMessage> removeMessageListener(noinline callback: (T) -> Unit){
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
        connectionStatusListeners.forEach{
            it.onFailedAttempt(attempt)
        }
    }

    override fun read(message: ApplicationMessage, mid: Int) {
        println("received message of type ${message.type}")
        for (callback: (ApplicationMessage) -> Unit in messageListeners.getOrDefault(message.javaClass, listOf<(ApplicationMessage) -> Unit>())) {
            println("invoking :)")
            callback.invoke(message)
        }


        for (callback: (ApplicationMessage) -> Unit in responseListeners.getOrDefault(mid, listOf<(ApplicationMessage) -> Unit>())) {
            println("invoking :)")
            callback.invoke(message)
        }
        responseListeners[mid]?.clear()
    }

    fun send(message: ApplicationMessage, callback: ((ApplicationMessage) -> Unit)? = null, desiredMessageId: Int = 0) {
        val json = message.toJson()


        callback?.let {
            addResponseListener(if (desiredMessageId != 0) desiredMessageId else messageId, callback)
        }


        tcpLayer?.write("${SimpleMessageReceiver.START_CHAR}${json.length}${SimpleMessageReceiver.SEPARATOR}${message.type}${SimpleMessageReceiver.SEPARATOR}${messageId}${SimpleMessageReceiver.SEPARATOR}$json")


        messageId++

        if (messageId > MESSAGE_ID_CEILING) {
            messageId = MESSAGE_STARTING_ID
        }
    }
}