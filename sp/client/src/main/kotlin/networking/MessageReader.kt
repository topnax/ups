package networking

interface MessageReader {
    fun read(message: Message): Unit
}