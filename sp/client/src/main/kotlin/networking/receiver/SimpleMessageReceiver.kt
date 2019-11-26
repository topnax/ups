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

    private var buffer = mutableListOf<Byte>()

    private var currentType: Int = 0
    private var currentLength: Int = 0
    private var currentID: Int = 0

    private var validHeader = false

    override fun receive(bytes: ByteArray, length: Int): List<ByteArray> {
        val message = String(bytes)

        val messagesBytes = mutableListOf<ByteArray>()
        var prevChar: Char? = null

        var lastGroupStart = 0
        logger.debug { "Received bytes of content '$message'" }

        message.forEachIndexed { i, char ->
            if (char == START_CHAR && (prevChar == null || prevChar != '\\') && lastGroupStart != i) {
                messagesBytes.add(bytes.copyOfRange(lastGroupStart, i))
                lastGroupStart = i
            }
            prevChar = char
        }

        if (lastGroupStart <= length && length <= bytes.size) {
            messagesBytes.add(bytes.copyOfRange(lastGroupStart, length))
        }

        messagesBytes.forEach {
            receiveMessage(it)
        }
        return messagesBytes
    }

    private fun receiveMessage(bytes: ByteArray) {
        val strMessage = String(bytes)
        logger.info { "Receiving message $strMessage" }
        if (!validHeader && strMessage[0] == START_CHAR && (buffer.isEmpty() || String(buffer.toByteArray())[buffer.size - 1] != '\\')) {
            val parts = strMessage.substring(1 until strMessage.length).split(SEPARATOR)
            if (parts.size == 4 && parts[0].isInt() && parts[1].isInt() && parts[2].isInt()) {
                validHeader = true
                currentLength = parts[0].toInt()
                currentType = parts[1].toInt()
                currentID = parts[2].toInt()
                bytes.copyOfRange(strMessage.indexOfNth(SEPARATOR, 3) + 1, bytes.size).forEach { buffer.add(it) }
                checkBuffer()
            } else {
                logger.error { "Receiver message '$strMessage' could not be parsed, invalid header." }
                validHeader = false
            }
        } else if (validHeader) {
            bytes.forEach { buffer.add(it) }
            checkBuffer()
        }
    }

    private fun checkBuffer() {
        if (validHeader && currentLength <= buffer.size) {
            messageReader.read(Message(currentLength, currentType, String(buffer.toByteArray()), currentID))
            currentLength = 0
            buffer.clear()
            validHeader = false
        }
    }
}

fun String.indexOfNth(char: Char, n: Int): Int {
    if (n < 1) {
        return -1
    }
    var lastIndex = 0
    var found = 0
    while (true) {
        lastIndex = this.indexOf(char, lastIndex)
        found++
        if (lastIndex != -1) {
            if (found == n) {
                return lastIndex
            }
            lastIndex++
        } else {
            return -1
        }
    }
}