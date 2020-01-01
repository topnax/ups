package networking.receiver

import mu.KotlinLogging
import networking.Network
import networking.reader.MessageReader

private val logger = KotlinLogging.logger {}

class SimpleMessageReceiver(messageReader: MessageReader, val testMode: Boolean = false) : MessageReceiver(messageReader) {

    companion object {
        const val START_CHAR = '$'
        const val SEPARATOR = '#'
        const val INVALID_BYTE_COUNT_CEILING = 5
    }

    private var state = 1
    private var closed = false
    private var type = 0
    private var length = 0
    private var messageId = 0
    private var invalidByteCount = 0
    private var headerBuffer = mutableListOf<Byte>()
    private var contentBuffer = mutableListOf<Byte>()

    override fun receive(bytes: ByteArray, length: Int) {
        val strMessage = String(bytes)
        logger.info { "Receiving message $strMessage" }

        // messages are parsed using an automaton
        bytes.forEachIndexed { index, byte ->
            if (index >= length) {
                return
            }
            when (state) {
                1 -> {
                    if (byte.toChar() == START_CHAR) {
                        state = 2
                        validByte()
                    } else {
                        invalidByte()
                    }
                }

                2 -> {
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            3
                        }
                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }
                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                3 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            3
                        }

                        byte.toChar() == SEPARATOR -> {
                            this.length = String(headerBuffer.toByteArray()).toInt()
                            validByte()
                            4
                        }

                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }

                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                4 -> {
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            5
                        }
                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }
                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                5 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            5
                        }

                        byte.toChar() == SEPARATOR -> {
                            type = String(headerBuffer.toByteArray()).toInt()
                            validByte()
                            6
                        }

                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }

                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                6 -> {
                    headerBuffer.clear()
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            7
                        }
                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }
                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                7 -> {
                    state = when {
                        byte.toChar().isDigit() -> {
                            headerBuffer.add(byte)
                            validByte()
                            7
                        }

                        byte.toChar() == SEPARATOR -> {
                            messageId = String(headerBuffer.toByteArray()).toInt()
                            contentBuffer.clear()
                            validByte()
                            8
                        }

                        byte.toChar() == START_CHAR -> {
                            invalidByte()
                            2
                        }

                        else -> {
                            invalidByte()
                            1
                        }
                    }
                }

                8 -> {
                    if (closed) {
                        return
                    }
                    state = if (!contentBuffer.isNextByteEscaped() && byte.toChar() == START_CHAR) {
                        invalidByte()
                        2
                    } else {
                        contentBuffer.add(byte)
                        when {
                            contentBuffer.size == this.length -> {
                                validByte()
                                messageReader.read(Message(this.length, type, String(contentBuffer.toByteArray()), messageId))
                                1
                            }
                            contentBuffer.size > this.length -> {
                                invalidByte()
                                8
                            }
                            else -> {
                                validByte()
                                8
                            }
                        }
                    }
                }
            }
        }
    }

    private fun invalidByte() {
        if (!testMode) {
            if (!closed) {
                invalidByteCount++
                logger.warn { "Received an invalid byte, invalidByteCount=$invalidByteCount" }
                if (invalidByteCount >= INVALID_BYTE_COUNT_CEILING) {
                    closed = true
                    logger.error { "Received invalidByteCount=$invalidByteCount, closing the connection." }
                    Network.getInstance().stop()
                }
            }
        }
    }

    private fun validByte() {
        invalidByteCount = 0
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
