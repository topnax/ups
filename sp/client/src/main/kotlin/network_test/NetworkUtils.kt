package hello

import java.io.IOException
import java.nio.ByteBuffer
import java.nio.channels.SelectionKey
import java.nio.channels.Selector
import java.nio.channels.SocketChannel

class NetworkUtils {
    companion object {

        fun readFromKey(key: SelectionKey): String? {
            try {
                val sc = key.channel() as SocketChannel
                val bb = ByteBuffer.allocate(1024)
                val len = sc.read(bb)
                return when (len) {
                    0, -1 -> {
                        null
                    }
                    else -> {
                        String(bb.array(), 0, len).trim()
                    }
                }
            } catch (e : IOException) {
                e.printStackTrace()
                return null
            }
        }

        fun writeToKey(msg: String?, key: SelectionKey) {
            msg?.let {
                val sChannel = key.channel() as SocketChannel
                val buffer = ByteBuffer.wrap(msg.toByteArray())
                sChannel.write(buffer)
            }
        }
    }
}