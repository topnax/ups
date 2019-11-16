package networking.messages

import com.beust.klaxon.*
import kotlin.reflect.KClass

class MessageTypeAdapter: TypeAdapter<ApplicationMessage> {
    override fun classFor(type: Any): KClass<out ApplicationMessage> = when(type as Int) {
        JoinLobbyMessage(0, "").type -> JoinLobbyMessage::class
        PlayerJoinedLobby(0, "").type -> PlayerJoinedLobby::class
        else -> throw IllegalArgumentException("Unknown type: $type")
    }
}

@TypeFor(field = "type", adapter = MessageTypeAdapter::class)
abstract class ApplicationMessage(@Json(ignored = true) val type: Int) {

    companion object {
        val renamer = object : FieldRenamer {
            override fun toJson(fieldName: String) = FieldRenamer.camelToUnderscores(fieldName)
            override fun fromJson(fieldName: String) = FieldRenamer.underscoreToCamel(fieldName)
        }

        inline fun <reified T> fromJson(json: String): T? where T : ApplicationMessage {
            return Klaxon().fieldRenamer(renamer).parse<T>(json)
        }


        fun fromJsonNew(json: String) : ApplicationMessage? {
            return Klaxon().fieldRenamer(renamer).parse<ApplicationMessage>(json)
        }

    }

    fun toJson(): String {
        return Klaxon().fieldRenamer(renamer).toJsonString(this)
    }
    inline fun <reified T> fromJson(json: String): T? where T : ApplicationMessage {
        return Klaxon().fieldRenamer(renamer).parse<T>(json)
    }

}

data class JoinLobbyMessage(val lobbyId: Int, val playerName: String) : ApplicationMessage(1)

data class PlayerJoinedLobby(val playerId: Int, val playerName: String) : ApplicationMessage(2)

fun main() {
    val json = """
        {
            "lobby_id": 10,
            "player_name": "Topinkos"
        }
    """

    val msg = ApplicationMessage.fromJsonNew(json)

//    Klaxon().converter()

    msg?.let {
        if (it is JoinLobbyMessage) {
            println(it.playerName)
            println(it.lobbyId)
        } else if (it is PlayerJoinedLobby) {
            println("pjl")
            println(it.playerName)
            println(it.lobbyId)
        }
    }
}
