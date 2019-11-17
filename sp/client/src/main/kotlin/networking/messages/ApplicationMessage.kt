package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon
import com.beust.klaxon.KlaxonException
import model.lobby.Lobby


abstract class ApplicationMessage(@Json(ignored = true) val type: Int) {

    companion object {
        val renamer = object : FieldRenamer {
            override fun toJson(fieldName: String) = FieldRenamer.camelToUnderscores(fieldName)
            override fun fromJson(fieldName: String) = FieldRenamer.underscoreToCamel(fieldName)
        }

        private inline fun <reified T> fromJson(json: String): T? where T : ApplicationMessage {
            return Klaxon().fieldRenamer(renamer).parse<T>(json)
        }

        fun fromJson(json: String, type: Int): ApplicationMessage? {
            return try {
                when (type) {
                    JoinLobbyMessage(0, "").type -> fromJson<JoinLobbyMessage>(json)
                    PlayerJoinedLobby(0, "").type -> fromJson<PlayerJoinedLobby>(json)
                    LobbiesListMessage(listOf()).type -> fromJson<LobbiesListMessage>(json)
                    else -> null
                }
            } catch (ex: KlaxonException) {
                null
            }
        }
    }

    fun toJson(): String {
        val json = Klaxon().fieldRenamer(renamer).toJsonString(this)
        println("parsed to '$json'")
        return json
    }
}

data class JoinLobbyMessage(val lobbyId: Int, val playerName: String) : ApplicationMessage(1)

data class PlayerJoinedLobby(val playerId: Int, val playerName: String) : ApplicationMessage(2)

class GetLobbiesMessage(val playerId: Int) : ApplicationMessage(3)

data class LobbiesListMessage(val lobbies: List<Lobby>) : ApplicationMessage(103)

fun main() {
    val msg: ApplicationMessage = JoinLobbyMessage(1,"topnax")
    println("dx" + msg.toJson())
    print(Klaxon().fieldRenamer(renamer = ApplicationMessage.renamer).toJsonString(JoinLobbyMessage(1, "topnax")))
}
