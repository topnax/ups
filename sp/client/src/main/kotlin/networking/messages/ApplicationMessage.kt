package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon
import com.beust.klaxon.KlaxonException


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
                    else -> null
                }
            } catch (ex: KlaxonException) {
                null
            }
        }
    }

    fun toJson(): String {
        return Klaxon().fieldRenamer(renamer).toJsonString(this)
    }
}

data class JoinLobbyMessage(val lobbyId: Int, val playerName: String) : ApplicationMessage(1)

data class PlayerJoinedLobby(val playerId: Int, val playerName: String) : ApplicationMessage(2)
