package networking.receiver

import mu.KotlinLogging
import networking.reader.MessageReader
import tornadofx.isInt

private val logger = KotlinLogging.logger {}

class FixedMessageReceiver(messageReader: MessageReader) : MessageReceiver(messageReader) {

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
        val messagesBytes = mutableListOf<ByteArray>()
        var lastMessageStart = 0
        var prevByte: Byte? = null

        bytes.forEachIndexed { index, byte ->
            if (lastMessageStart != index && byte == START_CHAR.toByte() && (prevByte == null || prevByte != '\\'.toByte())) {
                logger.debug { "MessageFEI added: '${String(bytes.copyOfRange(lastMessageStart, index))}'" }
                messagesBytes.add(bytes.copyOfRange(lastMessageStart, index))

                if (index < length) {
                    lastMessageStart = index
                }
            }
            prevByte = byte
        }

        if (lastMessageStart < length && length <= bytes.size) {
            logger.debug { "LastChance added: '${String(bytes.copyOfRange(lastMessageStart, length))}'" }
            messagesBytes.add(bytes.copyOfRange(lastMessageStart, length))
        }

//        messagesBytes.forEach {
//            receiveMessage(it)
//        }

        receiveMessage(bytes, length)
        return messagesBytes
    }

    private var state = 1
    private var autoType = 0
    private var autoLength = 0
    private var autoMessageId = 0
    private var automatonBuffer = mutableListOf<Byte>()
    private var automatonContentBuffer = mutableListOf<Byte>()


    private fun receiveMessage(bytes: ByteArray, length: Int) {
        val strMessage = String(bytes)
        logger.info { "Receiving message $strMessage" }

        bytes.forEachIndexed { index, byte ->
            if (index >= length) {
                return
            }
            val currentChar = byte.toChar()
            when (state) {
                1 -> {
                    if (byte.toChar() == START_CHAR) {
                        state = 2
                    }
                }

                2 -> {
                    automatonBuffer.clear()
                    state = if (byte.toChar().isDigit()) {
                        automatonBuffer.add(byte)
                        3
                    } else if (byte.toChar() == START_CHAR) {
                        2
                    } else {
                        1
                    }
                }

                3 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            automatonBuffer.add(byte)
                            3
                        }

                        byte.toChar() == SEPARATOR -> {
                            autoLength = String(automatonBuffer.toByteArray()).toInt()
                            4
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                4 -> {
                    automatonBuffer.clear()
                    state = if (byte.toChar().isDigit()) {
                        automatonBuffer.add(byte)
                        5
                    } else if (byte.toChar() == START_CHAR) {
                        2
                    } else {
                        1
                    }
                }

                5 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            automatonBuffer.add(byte)
                            5
                        }

                        byte.toChar() == SEPARATOR -> {
                            autoType = String(automatonBuffer.toByteArray()).toInt()
                            6
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                6 -> {
                    automatonBuffer.clear()
                    state = if (byte.toChar().isDigit()) {
                        automatonBuffer.add(byte)
                        7
                    } else if (byte.toChar() == START_CHAR) {
                        2
                    } else {
                        1
                    }
                }

                7 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            automatonBuffer.add(byte)
                            7
                        }

                        byte.toChar() == SEPARATOR -> {
                            autoMessageId = String(automatonBuffer.toByteArray()).toInt()
                            automatonContentBuffer.clear()
                            8
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                8 -> {
                    state = if (!automatonContentBuffer.isNextByteEscaped() && byte.toChar() == START_CHAR) {
                        2
                    } else {
                        automatonContentBuffer.add(byte)
                        if (automatonContentBuffer.size == autoLength) {
                            messageReader.read(Message(autoLength, autoType, String(automatonContentBuffer.toByteArray()), autoMessageId))
                            1
                        } else {
                            8
                        }
                    }
                }
            }
        }
    }

    private fun checkBuffer() {
        if (validHeader && currentLength <= buffer.size) {
            messageReader.read(Message(currentLength, currentType, String(buffer.toByteArray()).replace("\\", ""), currentID))
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
        if (lastIndex != -1) {
            found++
            if (found == n) {
                return lastIndex
            }
            lastIndex++
        } else {
            return -1
        }
    }
}


fun <Byte> MutableList<Byte>.isNextByteEscaped(): Boolean {
    var index = 0
    var escCount = 0
    while (this.size > 0 && this.size > index) {
        if ((this[this.size - index - 1] as kotlin.Byte).toChar() == '\\') {
            escCount++
        } else {
            break
        }
        index++
    }
    return escCount % 2 == 1
}
