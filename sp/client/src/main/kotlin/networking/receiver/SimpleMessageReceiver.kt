package networking.receiver

import mu.KotlinLogging
import networking.reader.MessageReader
import tornadofx.isInt

private val logger = KotlinLogging.logger {}

class SimpleMessageReceiver(messageReader: MessageReader) : MessageReceiver(messageReader) {

    companion object {
        val START_CHAR = '$'

        val SEPARATOR = '#'
    }

    private var buffer: String = ""

    private var currentType: Int = 0
    private var currentLength: Int = 0
    private var currentID: Int = 0

    private var validHeader = false

    override fun receive(bytes: ByteArray, length: Int) {
        val message = String(bytes)
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
    }

    private fun receiveMessage(message: String) {
        val length = message.length
        if (message[0] == START_CHAR && (buffer.isEmpty() || buffer[buffer.length - 1] != '\\')) {
            val parts = message.substring(1 until length).split(SEPARATOR)
            if (parts.size == 4 && parts[0].isInt() && parts[1].isInt() && parts[2].isInt()) {
                validHeader = true
                currentLength = parts[0].toInt()
                currentType = parts[1].toInt()
                currentID = parts[2].toInt()
                buffer = parts[3]
                checkBuffer()
            } else {
                logger.error { "Receiver message '$message' could not be parsed, invalid header." }
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
            messageReader.read(Message(currentLength, currentType, buffer, currentID))
            currentLength = 0
            buffer = ""
            validHeader = false
        }
    }
}
