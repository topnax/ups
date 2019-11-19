package networking.reader

import com.beust.klaxon.KlaxonException
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.JoinLobbyMessage
import networking.receiver.Message

class SimpleMessageReader(private val reader: ApplicationMessageReader) : MessageReader {

    private val messageTypes = hashMapOf<Int, Class<out ApplicationMessage>>()
    private val messageTypesX = hashMapOf<Int, (String) -> ApplicationMessage>()

    override fun read(message: Message) {
        println("Message read ${message.type} of content '${message.content}'")
        val am = ApplicationMessage.fromJson(message.content, message.type)
        am?.let {
            reader.read(am, message.id)
        }
    }
}