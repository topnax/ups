package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon
import com.beust.klaxon.KlaxonException
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
        const val LETTER_PLACED_MESSAGE_TYPE = 10
        const val LETTER_REMOVED_MESSAGE_TYPE = 11
        const val FINISH_ROUND_MESSAGE_TYPE = 12
        const val APPROVE_WORDS_MESSAGE_TYPE = 13
        const val DECLINE_WORDS_MESSAGE_TYPE = 14
        const val KEEP_ALIVE_MESSAGE_TYPE = 15
        const val LEAVE_GAME_MESSAGE_TYPE = 16

        const val GET_LOBBIES_RESPONSE_TYPE = 101
        const val LOBBY_UPDATED_RESPONSE_TYPE = 103
        const val LOBBY_DESTROYED_RESPONSE_TYPE = 105
        const val LOBBY_JOINED_RESPONSE_TYPE = 106
        const val USER_AUTHENTICATED_RESPONSE_TYPE = 107
        const val LOBBY_STARTED_RESPONSE_TYPE = 108
        const val GAME_STARTED_RESPONSE_TYPE = 109
        const val TILE_UPDATED_RESPONSE_TYPE = 110
        const val TILES_UPDATED_RESPONSE_TYPE = 111
        const val PLAYER_FINISHED_ROUND_RESPONSE_TYPE = 112
        const val PLAYER_ACCEPTED_ROUND_RESPONSE_TYPE = 113
        const val NEW_ROUND_RESPONSE_TYPE = 114
        const val YOUR_NEW_ROUND_RESPONSE_TYPE = 115
        const val PLAYER_DECLINED_WORDS_RESPONSE_TYPE = 116
        const val GAME_ENDED_RESPONSE_TYPE = 117
        const val ACCEPT_RESULTED_IN_NEW_ROUND_RESPONSE_TYPE = 118
        const val PLAYER_DISCONNECTED_RESPONSE_TYPE = 119
        const val GAME_STATE_REGENERATION_RESPONSE_TYPE = 120
        const val KEEP_ALIVE_RESPONSE_TYPE = 121
        const val USER_STATE_REGENERATION_RESPONSE_TYPE = 122
        const val FINISH_RESULTED_IN_NEXT_ROUND_TYPE = 123

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
                        LOBBY_STARTED_RESPONSE_TYPE -> fromJson<LobbyStartedResponse>(json)
                        GAME_STARTED_RESPONSE_TYPE -> fromJson<GameStartedResponse>(json)
                        TILE_UPDATED_RESPONSE_TYPE -> fromJson<TileUpdatedResponse>(json)
                        TILES_UPDATED_RESPONSE_TYPE -> fromJson<TilesUpdatedResponse>(json)
                        PLAYER_FINISHED_ROUND_RESPONSE_TYPE -> fromJson<RoundFinishedResponse>(json)
                        PLAYER_ACCEPTED_ROUND_RESPONSE_TYPE -> fromJson<PlayerAcceptedRoundResponse>(json)
                        NEW_ROUND_RESPONSE_TYPE -> fromJson<NewRoundResponse>(json)
                        YOUR_NEW_ROUND_RESPONSE_TYPE -> fromJson<YourNewRoundResponse>(json)
                        PLAYER_DECLINED_WORDS_RESPONSE_TYPE -> fromJson<PlayerDeclinedWordsResponse>(json)
                        GAME_ENDED_RESPONSE_TYPE -> fromJson<GameEndedResponse>(json)
                        ACCEPT_RESULTED_IN_NEW_ROUND_RESPONSE_TYPE -> fromJson<AcceptResultedInNewRound>(json)
                        PLAYER_DISCONNECTED_RESPONSE_TYPE -> fromJson<PlayerConnectionChangedResponse>(json)
                        GAME_STATE_REGENERATION_RESPONSE_TYPE -> fromJson<GameStateRegenerationResponse>(json)
                        KEEP_ALIVE_RESPONSE_TYPE -> fromJson<KeepAliveResponse>(json)
                        USER_STATE_REGENERATION_RESPONSE_TYPE -> fromJson<UserStateRegenerationResponse>(json)
                        FINISH_RESULTED_IN_NEXT_ROUND_TYPE -> fromJson<FinishResultedInNextRound>(json)
                        else -> null
                    }
                }
            } catch (ex: KlaxonException) {
                logger.error { "Failed to parse message of type $type from '$json', because '${ex.message}'" }
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
