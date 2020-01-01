package networking.messages

import model.game.Letter


data class ErrorResponseMessage(val content: String) : ApplicationMessage(ERROR_RESPONSE_TYPE)

class CreateLobbyMessage() : EmptyMessage(CREATE_LOBBY_MESSAGE_TYPE)

class GetLobbiesMessage : EmptyMessage(GET_LOBBIES_MESSAGE_TYPE)

data class JoinLobbyMessage(val lobbyId: Int) : ApplicationMessage(JOIN_LOBBY_MESSAGE_TYPE)

class LeaveLobbyMessage : EmptyMessage(LEAVE_LOBBY_MESSAGE_TYPE)

data class PlayerReadyToggleMessage(val ready: Boolean) : ApplicationMessage(PLAYER_READY_TOGGLE_MESSAGE_TYPE)

data class UserAuthenticationMessage(val name: String, val reconnecting: Boolean = false) : ApplicationMessage(USER_AUTHENTICATION_MESSAGE_TYPE)

class LetterPlacedMessage(val letter: Letter, val column: Int, val row: Int) : ApplicationMessage(LETTER_PLACED_MESSAGE_TYPE)

class LetterRemovedMessage(val column: Int, val row: Int) : ApplicationMessage(LETTER_REMOVED_MESSAGE_TYPE)

class UserLeavingMessage() : EmptyMessage(USER_LEAVING_MESSAGE_TYPE)

class StartLobbyMessage() : EmptyMessage(START_LOBBY_MESSAGE_TYPE)

class FinishRoundMessage() : EmptyMessage(FINISH_ROUND_MESSAGE_TYPE)

class ApproveWordsMessage() : EmptyMessage(APPROVE_WORDS_MESSAGE_TYPE)

class DeclineWordsMessage() : EmptyMessage(DECLINE_WORDS_MESSAGE_TYPE)

class KeepAliveMessage() : EmptyMessage(KEEP_ALIVE_MESSAGE_TYPE)

class LeaveGame() : EmptyMessage(LEAVE_GAME_MESSAGE_TYPE)
