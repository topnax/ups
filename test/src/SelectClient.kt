package hello

import java.io.IOException
import java.net.InetAddress
import java.net.InetSocketAddress
import java.nio.channels.SelectionKey
import java.nio.channels.Selector
import java.nio.channels.SocketChannel
import java.util.*
import kotlin.system.exitProcess


class SelectClient(val port: Int = 10000, val hostname: String = "localhost", val userName: String) {

    fun start() {
        val address = InetSocketAddress(InetAddress.getByName(hostname), port)
        val selector = Selector.open()
        val channel = SocketChannel.open()
        channel.configureBlocking(false)
        channel.connect(address)

        val operations = (SelectionKey.OP_CONNECT or SelectionKey.OP_READ
                or SelectionKey.OP_WRITE)

       channel.register(selector, operations)

        while (true) {
            if (selector.select() > 0) {
                val doneStatus = processReadySet(selector.selectedKeys())
                if (doneStatus) {
                    break
                }
            }
        }
        channel.close()
    }

    private fun listenForUserInput(key: SelectionKey) {
        Thread {
            val sc = Scanner(System.`in`)
            while (true) {
                val message = sc.nextLine()
                NetworkUtils.writeToKey(message, key)
            }
        }.start()
    }

    @Throws(Exception::class)
    fun processReadySet(readySet: MutableSet<*>): Boolean {
        val iterator = readySet.iterator()
        while (iterator.hasNext()) {
            val key = iterator.next() as SelectionKey
            iterator.remove()
            if (key.isConnectable) {
                if (!processConnect(key)) {
                    return true
                }
                NetworkUtils.writeToKey(userName, key)
                listenForUserInput(key)
            }
            if (key.isReadable) {
                val value = NetworkUtils.readFromKey(key)
                if (value == null) {
                    println("Server finished :)")
                    exitProcess(0)
                } else {
                    println(value)
                }
            }
        }
        return false
    }

    fun processConnect(key: SelectionKey): Boolean {
        val sc = key.channel() as SocketChannel
        try {
            while (sc.isConnectionPending) {
                sc.finishConnect()
            }
        } catch (e: IOException) {
            key.cancel()
            e.printStackTrace()
            return false
        }
        return true
    }
}
