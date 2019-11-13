package networking

import tornadofx.isInt
import tornadofx.launch
import java.util.logging.Level
import java.util.logging.Logger

class SimpleMessageReceiver(messageReader: MessageReader) : MessageReceiver(messageReader) {

    val LOG = Logger.getLogger(this.javaClass.name)

    companion object {
        val START_CHAR = '$'
        val SEPARATOR = '#'
    }

    private var buffer: String = ""

    override fun receive(bytes: ByteArray, length: Int) {
        val message = String(bytes)
        // if message[0] == START_CHAR && (len(s.buffers[UID].buffer) <= 0 || (len(buffer.buffer) > 0 && buffer.buffer[len(buffer.buffer)-1] != '\\')) {
        val messages = mutableListOf<String>()
        var prevChar: Char? = null

        var lastGroupStart = 0

        message.forEachIndexed { i, char ->
                if (char == '$' && (prevChar == null || prevChar != '\\') && lastGroupStart != i) {
                    messages.add(message.substring(lastGroupStart, i))
                    lastGroupStart = i
                }
        }

        messages.add(message.substring(lastGroupStart, message.length))

        print(messages)


//        println()me
//
//        LOG.level = Level.FINE
//        if (message[0] == START_CHAR && (buffer.isEmpty() || buffer[buffer.length - 1] != '\\')) {
//            println(message.substring(1 until length).split(SEPARATOR))
//            val parts = message.substring(1 until length).split(SEPARATOR)
//            if (parts[0].isInt() && parts[1].isInt()) {
//                println(parts[2].substring(0 until parts[0].toInt()))
//                println(message.indexOf(SEPARATOR, message.indexOf(SEPARATOR) + 1))
//                println("- " + message.substring(message.indexOf(SEPARATOR, message.indexOf(SEPARATOR) + 1) + parts[0].toInt() + 1, length))
//            } else {
//                LOG.severe("Receiver message '$message' could not be parsed, invalid header.")
//            }
//        }
    }
}

fun main(args: Array<String>) {
    val smr = SimpleMessageReceiver(object : MessageReader {
        override fun read(message: Message) {

        }
    }

    )

//    val message = "$6#1#{fele}"
    val message = "$6#1#{fele}$4#1#{fe}$"
    smr.receive(message.toByteArray(), message.length)

}