package networking.reader

import mu.KotlinLogging
import networking.Network
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.receiver.Message

private val logger = KotlinLogging.logger {}

class SimpleMessageReader(private val reader: ApplicationMessageReader) : MessageReader {

    companion object {
        const val CONSECUTIVE_INVALID_MESSAGE_LIMIT = 5
    }

    var consecutiveInvalidMessages = 0

    override fun read(message: Message) {
        logger.info { "Reading message of type ${message.type} and content '${message.content}'" }
        val am = ApplicationMessage.fromJson(message.content, message.type)
        am?.let {
            logger.debug { "About to read message of type ${am.javaClass.simpleName}" }
            reader.read(am, message.id)
            validMessageReceived()
        } ?: run {
            logger.error { "Could not read an ApplicationMessage because it could not be parsed" }
            invalidMessageReceived()
        }
    }

    private fun invalidMessageReceived() {
        consecutiveInvalidMessages++
        logger.warn { "Invalid message received. Consecutive invalid messages received=$consecutiveInvalidMessages" }
        if (consecutiveInvalidMessages >= CONSECUTIVE_INVALID_MESSAGE_LIMIT) {
            logger.error { "Consecutive invalid message count reached the ceiling of $CONSECUTIVE_INVALID_MESSAGE_LIMIT" }
            Network.getInstance().stop()
            consecutiveInvalidMessages = 0
        }
    }

    private fun validMessageReceived() {
        if (consecutiveInvalidMessages > 0) {
            logger.info { "Valid message received. Consecutive invalid message counter has been reset." }
        }
        consecutiveInvalidMessages = 0
    }
}