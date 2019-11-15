package networking.reader

import networking.receiver.Message

interface MessageReader {
    fun read(message: Message)
}