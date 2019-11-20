package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon
import com.beust.klaxon.KlaxonException
import model.lobby.Lobby
import model.lobby.Player


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
                if (type in 401..499) {
                    fromJson<ErrorResponseMessage>(json)
                } else {
                    when (type) {
                        JoinLobbyMessage(0, "").type -> fromJson<JoinLobbyMessage>(json)
                        PlayerJoinedLobby(0, "").type -> fromJson<PlayerJoinedLobby>(json)
                        GetLobbiesMessage().type -> fromJson<GetLobbiesMessage>(json)
                        SuccessResponseMessage("").type -> fromJson<SuccessResponseMessage>(json)
                        GetLobbiesResponse(mutableListOf()).type -> fromJson<GetLobbiesResponse>(json)
                        LobbyJoinedMessage(Lobby(listOf(), 0, Player("", 0))).type -> fromJson<LobbyJoinedMessage>(json)
                        LobbyDestroyedResponse().type -> fromJson<LobbyDestroyedResponse>(json)
                        else -> null
                    }
                }
            } catch (ex: KlaxonException) {
                println("json parse error!")
                println(json)

                println(ex)
                null
            }
        }
    }

    open fun toJson(): String {
        val json = Klaxon().fieldRenamer(renamer).toJsonString(this)
        println("parsed to '$json'")
        return json
    }
}

open class EmptyMessage(messageType: Int) : ApplicationMessage(messageType) {
    override fun toJson(): String {
        return "{}"
    }
}

data class PlayerJoinedLobby(val playerId: Int, val playerName: String) : ApplicationMessage(102)

data class CreateLobbyMessage(val playerName: String) : ApplicationMessage(2)

class GetLobbiesMessage : EmptyMessage(3)

data class JoinLobbyMessage(val lobbyId: Int, val playerName: String) : ApplicationMessage(4)

class LeaveLobbyMessage : EmptyMessage(5)

data class GetLobbiesResponse(val lobbies: MutableList<Lobby>): ApplicationMessage(101)

data class SuccessResponseMessage(val content: String) : ApplicationMessage(701)

data class ErrorResponseMessage(val content: String) : ApplicationMessage(101)

data class LobbyJoinedMessage(val lobby: Lobby) : ApplicationMessage(103)

class LobbyDestroyedResponse : EmptyMessage(105)

//class GetLobbiesMessage(val playerId: Int) : ApplicationMessage(102)

fun main() {
    try {
        val lobby = Klaxon().fieldRenamer(ApplicationMessage.renamer).parse<LobbyDestroyedResponse>("{}")
        println("Lobby type is ${lobby?.type}")
    } catch (ex: KlaxonException) {
        ex.printStackTrace()
    }
}
