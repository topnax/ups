package networking.receiver

import networking.reader.MessageReader

abstract class MessageReceiver(val messageReader: MessageReader) {
    abstract fun receive(bytes: ByteArray, length: Int) : List<ByteArray>
}