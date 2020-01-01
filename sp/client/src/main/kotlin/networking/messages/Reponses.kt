package networking.messages

import model.User
import model.game.Letter
import model.game.Tile
import model.lobby.Lobby
import model.lobby.Player

data class SuccessResponseMessage(val content: String) : ApplicationMessage(SUCCESS_RESPONSE_TYPE)

data class GetLobbiesResponse(val lobbies: MutableList<Lobby>) : ApplicationMessage(GET_LOBBIES_RESPONSE_TYPE)

data class LobbyUpdatedResponse(val lobby: Lobby) : ApplicationMessage(LOBBY_UPDATED_RESPONSE_TYPE)

data class LobbyJoinedResponse(val player: Player, val lobby: Lobby) : ApplicationMessage(LOBBY_JOINED_RESPONSE_TYPE)

class LobbyDestroyedResponse : EmptyMessage(LOBBY_DESTROYED_RESPONSE_TYPE)

data class UserAuthenticatedResponse(val user: User) : ApplicationMessage(USER_AUTHENTICATED_RESPONSE_TYPE)

class LobbyStartedResponse() : EmptyMessage(LOBBY_STARTED_RESPONSE_TYPE)

class GameStartedResponse(val players: List<Player>, val letters: List<Letter>, val activePlayerId: Int) : ApplicationMessage(GAME_STARTED_RESPONSE_TYPE)

class TileUpdatedResponse(val tile: Tile) : ApplicationMessage(TILE_UPDATED_RESPONSE_TYPE)

class TilesUpdatedResponse(val tiles: List<Tile>, val currentPlayerPoints: Int, val currentPlayerTotalPoints: Int) : ApplicationMessage(TILES_UPDATED_RESPONSE_TYPE)

class RoundFinishedResponse() : EmptyMessage(PLAYER_FINISHED_ROUND_RESPONSE_TYPE)

class PlayerAcceptedRoundResponse(val playerId: Int) : ApplicationMessage(PLAYER_ACCEPTED_ROUND_RESPONSE_TYPE)

class NewRoundResponse(val activePlayerId: Int) : ApplicationMessage(NEW_ROUND_RESPONSE_TYPE)

class YourNewRoundResponse(val letters: List<Letter>) : ApplicationMessage(YOUR_NEW_ROUND_RESPONSE_TYPE)

class PlayerDeclinedWordsResponse(val playerId: Int, val playerName: String) : ApplicationMessage(PLAYER_DECLINED_WORDS_RESPONSE_TYPE)

class GameEndedResponse(val playerPoints: Map<String, Player>) : ApplicationMessage(GAME_ENDED_RESPONSE_TYPE)

class AcceptResultedInNewRound() : EmptyMessage(ACCEPT_RESULTED_IN_NEW_ROUND_RESPONSE_TYPE)

class PlayerConnectionChangedResponse(val playerId: Int, val disconnected: Boolean) : ApplicationMessage(PLAYER_DISCONNECTED_RESPONSE_TYPE)

class GameStateRegenerationResponse(val user: User, val players: List<Player>, val tiles: List<Tile>, val activePlayerId: Int, val playerPoints: Map<String, Player>, val currentPlayerPoints: Int, val roundFinished: Boolean, val playerIdsThatAccepted: List<Int>, val letters: List<Letter>) : ApplicationMessage(GAME_STATE_REGENERATION_RESPONSE_TYPE)

class KeepAliveResponse() : EmptyMessage(KEEP_ALIVE_RESPONSE_TYPE)

class UserStateRegenerationResponse(var state: Int, val user: User? = null) : ApplicationMessage(USER_STATE_REGENERATION_RESPONSE_TYPE) {
    companion object {
        const val SERVER_RESTARTED = 0
        const val SERVER_RESTARTED_NAME_TAKEN = 1
        const val MOVED_TO_LOBBY_SCREEN = 2
        const val NOTHING = 3
    }

    init {
        if (state !in (SERVER_RESTARTED)..(NOTHING)) {
            state = SERVER_RESTARTED
        }
    }
}

class FinishResultedInNextRound() : EmptyMessage(FINISH_RESULTED_IN_NEXT_ROUND_TYPE)
