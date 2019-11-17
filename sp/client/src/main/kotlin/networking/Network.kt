package networking

import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
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

    private var tcpLayer: TCPLayer? = null

    private var messageListeners = mutableMapOf<Class<out ApplicationMessage>, MutableList<(ApplicationMessage) -> Unit>>()

    fun connectTo(hostname: String, port: Int) {
        tcpLayer = TCPLayer(port, hostname, SimpleMessageReceiver(SimpleMessageReader(this)), this)
    }

    fun addMessageListener(messageClazz: Class<out ApplicationMessage>, callback: (ApplicationMessage) -> Unit) {
        messageListeners.putIfAbsent(messageClazz, mutableListOf<(ApplicationMessage) -> Unit>())
        messageListeners[messageClazz]?.add(callback)
    }

    override fun onConnected() {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun onUnreachable() {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun onFailedAttempt(attempt: Int) {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }

    override fun read(message: ApplicationMessage) {
        TODO("not implemented") //To change body of created functions use File | Settings | File Templates.
    }
}