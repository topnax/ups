package networking

import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.GetLobbiesMessage
import networking.reader.SimpleMessageReader
import networking.receiver.MessageReceiver
import networking.receiver.SimpleMessageReceiver

class Network : ConnectionStatusListener, ApplicationMessageReader {

    companion object {

        private var network: Network = Network()

        @Synchronized
        fun getInstance(): Network {
            return network
        }
    }

    var tcpLayer: TCPLayer? = null

    private var messageListeners = mutableMapOf<Class<out ApplicationMessage>, MutableList<(ApplicationMessage) -> Unit>>()

    val connectionStatusListeners = mutableListOf<ConnectionStatusListener>()

    fun connectTo(hostname: String, port: Int) {
        tcpLayer?.close()
        tcpLayer = TCPLayer(port, hostname, SimpleMessageReceiver(SimpleMessageReader(this)), this)
        tcpLayer?.start()
    }

    fun addMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        messageListeners.putIfAbsent(messageClazz, mutableListOf())
        messageListeners[messageClazz]?.add(callback)
    }

    fun removeMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        messageListeners[messageClazz]?.remove(callback)
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

    override fun read(message: ApplicationMessage) {
        println("received message of type ${message.type}")
        for (callback: (ApplicationMessage) -> Unit in messageListeners.getOrDefault(message.javaClass, listOf<(ApplicationMessage) -> Unit>())) {
            println("invoking :)")
            callback.invoke(message)
        }
    }

    fun send(message: ApplicationMessage) {
        val json = message.toJson()

        tcpLayer?.write("${SimpleMessageReceiver.START_CHAR}${json.length}${SimpleMessageReceiver.SEPARATOR}${message.type}${SimpleMessageReceiver.SEPARATOR}$json")
    }
}