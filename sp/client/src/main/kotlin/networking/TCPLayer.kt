package networking

import networking.receiver.MessageReceiver
import java.io.*
import java.net.ConnectException
import java.net.InetAddress
import java.net.Socket
import java.net.SocketException

class TCPLayer(private val port: Int = 10000, val hostname: String = "localhost", val messageReceiver: MessageReceiver, val connectionStatusListener: ConnectionStatusListener) : Thread() {

    companion object {
        val NUMBER_OF_ATTEMPTS = 4
        val DELAY_BETWEEN_ATTEMPTS = 2000L
    }

    var socket: Socket? = null

    var output: OutputStream? = null

    var input: InputStream? = null

    private var run = true

    override fun run() {

        println("Client opening a socket at $hostname using port $port")
        for (i in 0..NUMBER_OF_ATTEMPTS) {
            try {
                socket = Socket(InetAddress.getByName(hostname), port)
                break
            } catch (exception: ConnectException) {
                connectionStatusListener.onFailedAttempt(i + 1)
                if (i != NUMBER_OF_ATTEMPTS) {
                    sleep(DELAY_BETWEEN_ATTEMPTS)
                }
            } catch (exception: Exception) {
                exception.printStackTrace()
                connectionStatusListener.onUnreachable()
            }
        }

        socket?.let {
            println("Socket created @ ${it.inetAddress} with port ${it.port}")
            output = it.getOutputStream()
            println("output gathered")
            input = it.getInputStream()
            println("output gathered")

            connectionStatusListener.onConnected()

            val serverMessage = ByteArray(100)

            try {
                println("Writing to server")

                while (run) {
                    var len = input?.read(serverMessage)

                    if (len == null) {
                        len = 0
                    }

                    val message: String? = when (len) {
                        0, -1 -> {
                            null
                        }
                        else -> {
                            String(serverMessage, 0, len)
                        }
                    }

                    message?.let {
                        messageReceiver.receive(serverMessage, len)
                        println("from server: '${it}', len ${len}")
                    }
                }
                println("stopped")

            } catch (e: IOException) {
                e.printStackTrace()
            } catch (e: SocketException) {
                connectionStatusListener.onUnreachable()
                e.printStackTrace()
            } finally {
                println("finished")
                close()
            }
        } ?: run {
            connectionStatusListener.onUnreachable()
        }
    }

    fun write(content: String) {
        output?.write(content.toByteArray())
    }

    fun close() {
        run = false
        input?.close()
        output?.close()
    }
}