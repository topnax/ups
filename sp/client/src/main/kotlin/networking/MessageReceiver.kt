package networking

abstract class MessageReceiver(val messageReader: MessageReader) {
    abstract fun receive(bytes: ByteArray, length: Int): Unit
}