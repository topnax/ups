package networking.applicationmessagereader

import networking.messages.ApplicationMessage
import networking.messages.PlayerJoinedLobby

class KrisKrosMessageReader : ApplicationMessageReader {

    val messageHandlers = hashMapOf<Class<out ApplicationMessage>, (ApplicationMessage) -> Unit>()

    init {
        register(PlayerJoinedLobby::class.java, ::onPlayerJoinedLobby)
    }

    fun register(clazz: Class<out ApplicationMessage>, handler: (ApplicationMessage) -> Unit) {
        messageHandlers[clazz] = handler
    }

    override fun read(message: ApplicationMessage, mid: Int) {
        messageHandlers[message.javaClass]?.invoke(message)
    }

    fun onPlayerJoinedLobby(message: ApplicationMessage) {
        if (message is PlayerJoinedLobby) {
            println("Player ${message.playerName} joined the lobby")
        }
    }
}