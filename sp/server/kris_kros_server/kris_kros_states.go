package kris_kros_server

import (
	"ups/sp/server/protocol/messages"
)

const (
	INITIAL_STATE_ID                = 1  // a state of unauthorized users
	AUTHORIZED_STATE_ID             = 2  // state of users that have been authorized and can view available lobbies and create them
	LOBBY_JOINED_ID                 = 3  // state of players that have joined a lobby and can toggle their ready state
	LOBBY_JOINED_READY_ID           = 4  // state of players that are ready and have joined a lobby and can toggle their ready state
	LOBBY_CREATED_ID                = 5  // state of players that have created a lobby and can toggle their ready state
	LOBBY_CREATED_READY_ID          = 6  // state of players that are ready and have created a lobby and can toggle their ready state
	GAME_STARTED_STATE_ID           = 7  // state of players that have started a game
	PLAYERS_TURN_STATE_ID           = 8  // state of players whose turn it is
	PLAYER_WAITING_ID               = 9  // state of players whose turn it isn't
	PLAYER_FINISHED_ROUND_ID        = 10 // state of players that have finished their turn
	APPROVE_WORDS_STATE_ID          = 11 // state of players that should decide words validity
	WORDS_VALIDITY_DECIDED_STATE_ID = 12 // state of players that have already decided words validity
)

type State interface {
	Id() int
	Routes() map[int]int
}

// registers the possible states of the users
func (router *KrisKrosRouter) registerStates() {
	router.states = make(map[int]State)
	router.registerState(InitialState{})
	router.registerState(AuthorizedState{})
	router.registerState(LobbyJoinedState{})
	router.registerState(LobbyJoinedReadyState{})
	router.registerState(LobbyCreatedState{})
	router.registerState(LobbyCreatedReadyState{})
	router.registerState(GameStartedState{})
	router.registerState(PlayersTurnState{})
	router.registerState(PlayerWaitingState{})
	router.registerState(PlayerFinishedRoundState{})
	router.registerState(ApproveWordsState{})
	router.registerState(WordsValidityDecidedState{})
}

////////////////////////////////////////////

type InitialState struct{}

func (i InitialState) Id() int {
	return INITIAL_STATE_ID
}

func (i InitialState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.UserAuthenticationMessageType] = AuthorizedState{}.Id()
	m[messages.UserLeavingMessageType] = InitialState{}.Id()
	return m
}

////////////////////////////////////////////

type AuthorizedState struct{}

func (a AuthorizedState) Id() int {
	return AUTHORIZED_STATE_ID
}

func (a AuthorizedState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.GetLobbiesType] = AuthorizedState{}.Id()
	m[messages.JoinLobbyMessageType] = LobbyJoinedState{}.Id()
	m[messages.CreateLobbyType] = LobbyCreatedState{}.Id()
	return m
}

////////////////////////////////////////////

type LobbyJoinedState struct{}

func (a LobbyJoinedState) Id() int {
	return LOBBY_JOINED_ID
}

func (a LobbyJoinedState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveLobbyMessageType] = AuthorizedState{}.Id()
	m[messages.PlayerReadyMessageType] = LobbyJoinedReadyState{}.Id()
	return m
}

////////////////////////////////////////////

type LobbyJoinedReadyState struct{}

func (a LobbyJoinedReadyState) Id() int {
	return LOBBY_JOINED_READY_ID
}

func (a LobbyJoinedReadyState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.PlayerReadyMessageType] = LobbyJoinedState{}.Id()
	m[messages.LeaveLobbyMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type LobbyCreatedState struct{}

func (a LobbyCreatedState) Id() int {
	return LOBBY_CREATED_ID
}

func (a LobbyCreatedState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveLobbyMessageType] = AuthorizedState{}.Id()
	m[messages.PlayerReadyMessageType] = LobbyCreatedReadyState{}.Id()
	return m
}

////////////////////////////////////////////

type LobbyCreatedReadyState struct{}

func (a LobbyCreatedReadyState) Id() int {
	return LOBBY_CREATED_READY_ID
}

func (a LobbyCreatedReadyState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.PlayerReadyMessageType] = LobbyCreatedState{}.Id()
	m[messages.LeaveLobbyMessageType] = AuthorizedState{}.Id()
	m[messages.StartLobbyMessageType] = GameStartedState{}.Id()
	return m
}

////////////////////////////////////////////

type GameStartedState struct{}

func (a GameStartedState) Id() int {
	return GAME_STARTED_STATE_ID
}

func (a GameStartedState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveLobbyMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type PlayersTurnState struct{}

func (a PlayersTurnState) Id() int {
	return PLAYERS_TURN_STATE_ID
}

func (a PlayersTurnState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LetterPlacedMessageType] = PlayersTurnState{}.Id()
	m[messages.LetterRemovedMessageType] = PlayersTurnState{}.Id()
	m[messages.FinishRoundMessageType] = PlayerFinishedRoundState{}.Id()
	m[messages.LeaveGameMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type PlayerWaitingState struct{}

func (a PlayerWaitingState) Id() int {
	return PLAYER_WAITING_ID
}

func (a PlayerWaitingState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveGameMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type PlayerFinishedRoundState struct{}

func (a PlayerFinishedRoundState) Id() int {
	return PLAYER_FINISHED_ROUND_ID
}

func (a PlayerFinishedRoundState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveGameMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type ApproveWordsState struct{}

func (a ApproveWordsState) Id() int {
	return APPROVE_WORDS_STATE_ID
}

func (a ApproveWordsState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.ApproveWordsMessageType] = WordsValidityDecidedState{}.Id()
	m[messages.DeclineWordsMessageType] = WordsValidityDecidedState{}.Id()
	m[messages.LeaveGameMessageType] = AuthorizedState{}.Id()
	return m
}

////////////////////////////////////////////

type WordsValidityDecidedState struct{}

func (a WordsValidityDecidedState) Id() int {
	return WORDS_VALIDITY_DECIDED_STATE_ID
}

func (a WordsValidityDecidedState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.LeaveGameMessageType] = AuthorizedState{}.Id()
	return m
}
