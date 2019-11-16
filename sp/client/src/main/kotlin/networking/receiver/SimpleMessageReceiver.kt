package networking.receiver

import networking.reader.MessageReader
import tornadofx.isInt
import java.util.logging.Logger

class SimpleMessageReceiver(messageReader: MessageReader) : MessageReceiver(messageReader) {

    val LOG = Logger.getLogger(this.javaClass.name)
    companion object {

        val START_CHAR = '$'
        val SEPARATOR = '#'
    }
    private var buffer: String = ""

    private var currentLength: Int = 0
    private var currentType: Int = 0

    private var validHeader = false

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
            prevChar = char
        }

        messages.add(message.substring(lastGroupStart, message.length))

        messages.forEach {
            receiveMessage(it)
        }

//        println()me
//
//        LOG.level = Level.FINE

    }

    private fun receiveMessage(message: String) {
        val length = message.length
        if (message[0] == START_CHAR && (buffer.isEmpty() || buffer[buffer.length - 1] != '\\')) {
            println(message.substring(1 until length).split(SEPARATOR))
            val parts = message.substring(1 until length).split(SEPARATOR)
            if (parts.size == 3 && parts[0].isInt() && parts[1].isInt()) {
                validHeader = true
                currentLength = parts[0].toInt()
                currentType = parts[1].toInt()
                buffer = parts[2]
                checkBuffer()
            } else {
                LOG.severe("Receiver message '$message' could not be parsed, invalid header.")
                validHeader = false
            }
        } else {
            buffer += message
            checkBuffer()
        }
    }

    private fun checkBuffer() {
        if (validHeader && currentLength <= buffer.length) {
            buffer = buffer.substring(0 until currentLength)
            messageReader.read(Message(currentLength, currentType, buffer))
            println("success $currentType :) '$buffer'")
            currentLength = 0
            buffer = ""
            validHeader = false
        }
    }
}