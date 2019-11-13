package networking

import hello.server
import java.io.*
import java.net.ConnectException
import java.net.InetAddress
import java.net.Socket
import java.nio.ByteBuffer
import java.nio.charset.Charset

class TCPLayer(private val port: Int = 10000, val hostname: String = "localhost", val messageReceiver: MessageReceiver, val connectionStatusListener: ConnectionStatusListener) : Thread() {

    companion object {
        val NUMBER_OF_ATTEMPTS = 4
        val DELAY_BETWEEN_ATTEMPTS = 250L
    }

    var socket: Socket? = null

    lateinit var output: OutputStream

    lateinit var input: InputStream

    init {


    }

    override fun run() {

        println("Client opening a socket at $hostname using port $port")
        for (i in 0..NUMBER_OF_ATTEMPTS) {
            try {
                socket = Socket(InetAddress.getByName(hostname), port)
                break;
            } catch (exception: ConnectException) {
                connectionStatusListener.onFailedAttempt(i + 1)
                if (i != NUMBER_OF_ATTEMPTS) {
                    sleep(DELAY_BETWEEN_ATTEMPTS)
                }
            }
        }

        socket?.let {
            println("Socket created @ ${it.inetAddress} with port ${it.port}")
            output = it.getOutputStream()
            println("output gathered")
            input = it.getInputStream()
            println("output gathered")

            connectionStatusListener.onConnected()

            var serverMessage: ByteArray = ByteArray(100)

            try {
                println("Writing to server")

                while (true) {

                    val len = input.read(serverMessage)
                    val message: String? = when (len) {
                        0, -1 -> {
                            null
                        }
                        else -> {
                            String(serverMessage, 0, len)
                        }
                    }

                    message?.let {
                        messageReceiver.receive(serverMessage.toTypedArray(), len)
                        println("from server: '${it}', len ${len}")
                    }
                }
            } catch (e: IOException) {
                e.printStackTrace()
            } finally {
                println("finished")
                finish()
            }
        } ?: run {
            connectionStatusListener.onUnreachable()
        }
    }

    private fun finish() {
//        input.close()
        output.close()
    }
}