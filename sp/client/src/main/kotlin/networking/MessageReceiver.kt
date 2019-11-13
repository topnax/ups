package networking

abstract class MessageReceiver(val messageReader: MessageReader) {
    abstract fun receive(bytes: Array<Byte>, length: Int): Unit
}