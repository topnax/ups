package networking.reader

import mu.KotlinLogging
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.receiver.Message

private val logger = KotlinLogging.logger {}

class SimpleMessageReader(private val reader: ApplicationMessageReader) : MessageReader {
    override fun read(message: Message) {
        logger.info { "Reading message of type ${message.type} and content '${message.content}'" }
        val am = ApplicationMessage.fromJson(message.content, message.type)
        am?.let {
            logger.debug { "About to read message of type ${am.javaClass.simpleName}" }
            reader.read(am, message.id)
        } ?: run {
            logger.error { "Could not read an ApplicationMessage because it could not be parsed" }
        }
    }
}