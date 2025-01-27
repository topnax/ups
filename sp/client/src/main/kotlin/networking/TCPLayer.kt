package networking

import mu.KotlinLogging
import networking.receiver.MessageReceiver
import java.io.IOException
import java.io.InputStream
import java.io.OutputStream
import java.net.ConnectException
import java.net.InetAddress
import java.net.NoRouteToHostException
import java.net.Socket

private val logger = KotlinLogging.logger {}

class TCPLayer(private val port: Int = 10000, val hostname: String = "localhost", val messageReceiver: MessageReceiver, val connectionStatusListener: ConnectionStatusListener) : Thread() {

    companion object {
        val NUMBER_OF_ATTEMPTS = 4
        val DELAY_BETWEEN_ATTEMPTS = 3000L
    }

    var socket: Socket? = null

    var output: OutputStream? = null

    var input: InputStream? = null

    private var run = true

    override fun run() {
        connect()
    }

    fun connect() {
        logger.info { "Opening a socket at $hostname using port $port" }
        // limited number of attempts before onUnreachable event is raised
        for (i in 0..NUMBER_OF_ATTEMPTS) {
            try {
                socket = Socket(InetAddress.getByName(hostname), port)
                break
            } catch (exception: ConnectException) {
                connectionStatusListener.onFailedAttempt(i + 1)
                if (i != NUMBER_OF_ATTEMPTS) {
                    sleep(DELAY_BETWEEN_ATTEMPTS)
                }
            } catch (exception: NoRouteToHostException) {
                connectionStatusListener.onFailedAttempt(i + 1)
                if (i != NUMBER_OF_ATTEMPTS) {
                    sleep(DELAY_BETWEEN_ATTEMPTS)
                }
            } catch (exception: Exception) {
                logger.error(exception) { "Could not open a socket to the server" }
                connectionStatusListener.onUnreachable()
            }
        }

        // start a read loop
        socket?.let {
            logger.info { "Socket successfully created @ ${it.inetAddress} with port ${it.port}" }
            output = it.getOutputStream()
            input = it.getInputStream()

            connectionStatusListener.onConnected()

            var messageBytes = ByteArray(100)

            try {
                while (run) {
                    var len = input?.read(messageBytes)

                    if (len == null) {
                        len = 0
                    }

                    val message: String? = when (len) {
                        0, -1 -> {
                            null
                        }
                        else -> {
                            String(messageBytes, 0, len)
                        }
                    }

                    message?.let {
                        logger.info { "received from server str'$it' of length $len" }
                        messageReceiver.receive(messageBytes, len)
                    }
                    messageBytes = ByteArray(100)
                }
                logger.info { "tcp read loop stopped" }
            } catch (e: IOException) {
                logger.error(e) { "a socket has raised an exception" }
            } finally {
                close()
            }
        } ?: run {
            connectionStatusListener.onUnreachable()
        }
    }

    fun write(content: String) {
        try {
            logger.info { "Writing to server '$content'" }
            output?.write(content.toByteArray())
        } catch (ex: IOException) {
            logger.error(ex) { "an exception was thrown during writing to server" }
        }
    }

    fun close() {
        logger.info { "Closing TCP layer..." }
        run = false
        input?.close()
        output?.close()
    }
}