package networking

abstract class MessageReceiver(val messageReader: MessageReader) {
    abstract fun receive(UID: Int, bytes: Array<Byte>, length: Int): Unit
}