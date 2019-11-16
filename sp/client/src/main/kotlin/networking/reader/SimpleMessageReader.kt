package networking.reader

import networking.messages.ApplicationMessage
import networking.messages.JoinLobbyMessage
import networking.messages.PlayerJoinedLobby
import networking.receiver.Message

class SimpleMessageReader : MessageReader {

    private val messageTypes = hashMapOf<Int, Class<out ApplicationMessage>>()
    private val messageTypesX = hashMapOf<Int, (String) -> ApplicationMessage >()

    init {
        register(JoinLobbyMessage(0,""))
        register(PlayerJoinedLobby(0,""))
    }

    private fun register(message: ApplicationMessage) {
        messageTypes.put(message.type, message.javaClass)
        messageTypesX.put(message.type, {str: String -> return ApplicationMessage.fromJson<>()message.javaClass} )
    }

    override fun read(message: Message) {
        val messageType = messageTypes[message.type]
        val messageTypeX: ApplicationMessage? = messageTypesX[message.type]
        if (messageTypeX != null) {
//            ApplicationMessage.fromJson<messageType>(message.content)
            messageTypeX.fromJson<ApplicationMessage>()
        }
    }
}