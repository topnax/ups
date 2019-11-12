package hello

import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.InetAddress
import java.net.InetSocketAddress
import java.nio.ByteBuffer
import java.nio.channels.SelectionKey
import java.nio.channels.Selector
import java.nio.channels.ServerSocketChannel
import java.nio.channels.SocketChannel


class SelectServer(val port: Int = 10000, val hostname: String = "localhost") {

    private val clients = mutableSetOf<SocketChannel>()
    private val clientNames = mutableMapOf<SocketChannel, String>()

    fun start() {
        val address = InetAddress.getByName(hostname)
        val selector = Selector.open()

        val serverSocketChannel = ServerSocketChannel.open()

        serverSocketChannel.configureBlocking(false)
        serverSocketChannel.bind(InetSocketAddress(address, port))
        serverSocketChannel.register(selector, SelectionKey.OP_ACCEPT)

        while (true) {
            if (selector.select() > 0) {
                val selectedKeys = selector.selectedKeys()
                val listToRemove = mutableListOf<SelectionKey>()
                for (key in selectedKeys) {
                    if (key.isAcceptable) {
                        val sc = serverSocketChannel.accept()
                        sc?.let {
                            clients.add(sc)
                            it.configureBlocking(false)
                            it.register(selector, SelectionKey.OP_READ)
                            println(
                                "Connection Accepted: ${it.localAddress} \n ${it} ${it.socket().localPort}"
                            )
                        }
                    }
                    if (key.isReadable) {
                        val sc = key.channel() as SocketChannel
                        val readValue = NetworkUtils.readFromKey(key)
                        when (readValue) {
                            null -> exitClient(sc)
                            else -> {
                                var message: String?

                                if (clientNames[sc] == null) {
                                    clientNames[sc] = readValue
                                    message = "$readValue has joined the chat...\n"
                                } else {
                                    message = "[${clientNames[sc]}] - $readValue\n"
                                }

                                broadcast(message, sc)

                                listToRemove.add(key)
                            }
                        }
                    }
                }
                selectedKeys.clear()
            }
        }
    }

    private fun broadcast(message: String, exclude: SocketChannel? = null) {
        println("[BCAST]: '$message'")
        for (csc in clients) {
            if (exclude == null || exclude != csc) {
                val bb = ByteBuffer.wrap(message.toByteArray())
                csc.write(bb)
            }
        }
    }

    private fun exitClient(sc: SocketChannel) {
        val name = if (clientNames[sc] != null) {
            clientNames[sc]
        } else {
            "unknown client"
        }
        broadcast("[$name] has left the chat\n", sc)
        sc.close()
        clients.remove(sc)
        clientNames.remove(sc)
    }
}