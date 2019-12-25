package networking.receiver

import mu.KotlinLogging
import networking.reader.MessageReader

private val logger = KotlinLogging.logger {}

class FixedMessageReceiver(messageReader: MessageReader) : MessageReceiver(messageReader) {

    companion object {
        val START_CHAR = '$'
        val SEPARATOR = '#'
    }

    private var state = 1
    private var type = 0
    private var length = 0
    private var messageId = 0
    private var headerBuffer = mutableListOf<Byte>()
    private var contentBuffer = mutableListOf<Byte>()

    override fun receive(bytes: ByteArray, length: Int) {
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
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            3
                        }
                        byte.toChar() == START_CHAR -> 2
                        else -> 1
                    }
                }

                3 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            3
                        }

                        byte.toChar() == SEPARATOR -> {
                            this.length = String(headerBuffer.toByteArray()).toInt()
                            4
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                4 -> {
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            5
                        }
                        byte.toChar() == START_CHAR -> 2
                        else -> 1
                    }
                }

                5 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            5
                        }

                        byte.toChar() == SEPARATOR -> {
                            type = String(headerBuffer.toByteArray()).toInt()
                            6
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                6 -> {
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            7
                        }
                        byte.toChar() == START_CHAR -> 2
                        else -> 1
                    }
                }

                7 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            7
                        }

                        byte.toChar() == SEPARATOR -> {
                            messageId = String(headerBuffer.toByteArray()).toInt()
                            contentBuffer.clear()
                            8
                        }

                        byte.toChar() == START_CHAR -> 2

                        else -> {
                            1
                        }
                    }
                }

                8 -> {
                    state = if (!contentBuffer.isNextByteEscaped() && byte.toChar() == START_CHAR) {
                        2
                    } else {
                        contentBuffer.add(byte)
                        if (contentBuffer.size == this.length) {
                            messageReader.read(Message(this.length, type, String(contentBuffer.toByteArray()), messageId))
                            1
                        } else {
                            8
                        }
                    }
                }
            }
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
