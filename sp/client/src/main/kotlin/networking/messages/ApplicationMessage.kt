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
                if (type in 401..499) {
                    fromJson<ErrorResponseMessage>(json)
                } else {
                    when (type) {
                        JoinLobbyMessage(0, "").type -> fromJson<JoinLobbyMessage>(json)
                        PlayerJoinedLobby(0, "").type -> fromJson<PlayerJoinedLobby>(json)
                        GetLobbiesMessage(listOf()).type -> fromJson<GetLobbiesMessage>(json)
                        SuccessResponseMessage("").type -> fromJson<SuccessResponseMessage>(json)
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

    fun toJson(): String {
        val json = Klaxon().fieldRenamer(renamer).toJsonString(this)
        println("parsed to '$json'")
        return json
    }
}

data class PlayerJoinedLobby(val playerId: Int, val playerName: String) : ApplicationMessage(102)

data class CreateLobbyMessage(val clientName: String) : ApplicationMessage(2)

data class GetLobbiesMessage(val lobbies: List<Lobby>) : ApplicationMessage(3)

data class JoinLobbyMessage(val lobbyId: Int, val clientName: String) : ApplicationMessage(4)

data class SuccessResponseMessage(val content: String) : ApplicationMessage(701)

data class ErrorResponseMessage(val content: String) : ApplicationMessage(101)

//class GetLobbiesMessage(val playerId: Int) : ApplicationMessage(102)

fun main() {
    val msg: ApplicationMessage = JoinLobbyMessage(1,"topnax")
    println("dx" + msg.toJson())
    print(Klaxon().fieldRenamer(renamer = ApplicationMessage.renamer).toJsonString(JoinLobbyMessage(1, "topnax")))
}
