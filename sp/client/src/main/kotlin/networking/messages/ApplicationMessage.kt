package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon
import com.beust.klaxon.KlaxonException
import model.lobby.Lobby
import model.lobby.Player
import model.lobby.User
import mu.KotlinLogging

val logger = KotlinLogging.logger { }

abstract class ApplicationMessage(@Json(ignored = true) val type: Int) {

    companion object {

        const val CREATE_LOBBY_MESSAGE_TYPE = 2
        const val GET_LOBBIES_MESSAGE_TYPE = 3
        const val JOIN_LOBBY_MESSAGE_TYPE = 4
        const val LEAVE_LOBBY_MESSAGE_TYPE = 5
        const val PLAYER_READY_TOGGLE_MESSAGE_TYPE = 6
        const val USER_AUTHENTICATION_MESSAGE_TYPE = 7
        const val USER_LEAVING_MESSAGE_TYPE = 8
        const val START_LOBBY_MESSAGE_TYPE = 9

        const val GET_LOBBIES_RESPONSE_TYPE = 101
        const val LOBBY_UPDATED_RESPONSE_TYPE = 103
        const val LOBBY_DESTROYED_RESPONSE_TYPE = 105
        const val LOBBY_JOINED_RESPONSE_TYPE = 106
        const val USER_AUTHENTICATED_RESPONSE_TYPE = 107
        const val LOBBY_STARTED_MESSAGE_TYPE = 108

        const val ERROR_RESPONSE_TYPE = 401
        const val SUCCESS_RESPONSE_TYPE = 701

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
                    val error = fromJson<ErrorResponseMessage>(json)
                    error?.let { logger.error { "Received an error of type $type and content ${error.content}" } }
                    error
                } else {
                    when (type) {
                        SUCCESS_RESPONSE_TYPE -> fromJson<SuccessResponseMessage>(json)
                        GET_LOBBIES_RESPONSE_TYPE -> fromJson<GetLobbiesResponse>(json)
                        LOBBY_UPDATED_RESPONSE_TYPE -> fromJson<LobbyUpdatedResponse>(json)
                        LOBBY_DESTROYED_RESPONSE_TYPE -> fromJson<LobbyDestroyedResponse>(json)
                        LOBBY_JOINED_RESPONSE_TYPE -> fromJson<LobbyJoinedResponse>(json)
                        USER_AUTHENTICATED_RESPONSE_TYPE -> fromJson<UserAuthenticatedResponse>(json)
                        LOBBY_STARTED_MESSAGE_TYPE -> fromJson<LobbyStartedResponse>(json)
                        else -> null
                    }
                }
            } catch (ex: KlaxonException) {
                logger.error { "Failed to parse message of type $type from '$json', because ${ex.message}" }
                null
            }
        }
    }

    open fun toJson(): String {
        val json = Klaxon().fieldRenamer(renamer).toJsonString(this)
        logger.info { "${this.javaClass::class.java} parsed to '$json'" }
        return json
    }
}

open class EmptyMessage(messageType: Int) : ApplicationMessage(messageType) {
    override fun toJson(): String {
        return "{}"
    }
}

data class ErrorResponseMessage(val content: String) : ApplicationMessage(ERROR_RESPONSE_TYPE)

class CreateLobbyMessage() : EmptyMessage(CREATE_LOBBY_MESSAGE_TYPE)

class GetLobbiesMessage : EmptyMessage(GET_LOBBIES_MESSAGE_TYPE)

data class JoinLobbyMessage(val lobbyId: Int) : ApplicationMessage(JOIN_LOBBY_MESSAGE_TYPE)

class LeaveLobbyMessage : EmptyMessage(LEAVE_LOBBY_MESSAGE_TYPE)

data class PlayerReadyToggleMessage(val ready: Boolean) : ApplicationMessage(PLAYER_READY_TOGGLE_MESSAGE_TYPE)

data class GetLobbiesResponse(val lobbies: MutableList<Lobby>) : ApplicationMessage(GET_LOBBIES_RESPONSE_TYPE)

data class LobbyUpdatedResponse(val lobby: Lobby) : ApplicationMessage(LOBBY_UPDATED_RESPONSE_TYPE)

data class LobbyJoinedResponse(val player: Player, val lobby: Lobby) : ApplicationMessage(LOBBY_JOINED_RESPONSE_TYPE)

class LobbyDestroyedResponse : EmptyMessage(LOBBY_DESTROYED_RESPONSE_TYPE)

data class SuccessResponseMessage(val content: String) : ApplicationMessage(SUCCESS_RESPONSE_TYPE)

data class UserAuthenticationMessage(val name: String) : ApplicationMessage(USER_AUTHENTICATION_MESSAGE_TYPE)

data class UserAuthenticatedResponse(val user: User) : ApplicationMessage(USER_AUTHENTICATED_RESPONSE_TYPE)

class UserLeavingMessage() : EmptyMessage(USER_LEAVING_MESSAGE_TYPE)

class StartLobbyMessage() : EmptyMessage(START_LOBBY_MESSAGE_TYPE)

class LobbyStartedResponse() : EmptyMessage(LOBBY_STARTED_MESSAGE_TYPE)
