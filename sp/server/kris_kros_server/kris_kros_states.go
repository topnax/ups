package kris_kros_server

import (
	"ups/sp/server/protocol/messages"
)

const (
	INITIAL_STATE_ID         = 1
	AUTHORIZED_STATE_ID      = 2
	LOBBY_JOINED_ID          = 3
	LOBBY_JOINED_READY_ID    = 4
	LOBBY_CREATED_ID         = 5
	LOBBY_CREATED_READY_ID   = 6
	GAME_STARTED_STATE_ID    = 7
	PLAYERS_TURN_STATE_ID    = 8
	PLAYER_WAITING_ID        = 9
	PLAYER_FINISHED_ROUND_ID = 10
	APPROVE_WORDS_STATE_ID   = 11
	APPROVED_WORDS_STATE_ID  = 12
)

type State interface {
	Id() int
	Routes() map[int]int
}

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
}

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
	m[messages.FinishRoundMessageType] = FinishedRoundMessage{}.Id()
	return m
}

////////////////////////////////////////////

type PlayerWaitingState struct{}

func (a PlayerWaitingState) Id() int {
	return PLAYER_WAITING_ID
}

func (a PlayerWaitingState) Routes() map[int]int {
	m := make(map[int]int)
	return m
}

////////////////////////////////////////////

type FinishedRoundMessage struct{}

func (a FinishedRoundMessage) Id() int {
	return PLAYER_FINISHED_ROUND_ID
}

func (a FinishedRoundMessage) Routes() map[int]int {
	m := make(map[int]int)
	return m
}

////////////////////////////////////////////

type ApproveWordsState struct{}

func (a ApproveWordsState) Id() int {
	return APPROVE_WORDS_STATE_ID
}

func (a ApproveWordsState) Routes() map[int]int {
	m := make(map[int]int)
	m[messages.ApproveWordsMessageType] = PlayersTurnState{}.Id()
	m[messages.DeclineWordsMessageType] = PlayerWaitingState{}.Id()
	return m
}

////////////////////////////////////////////

type ApprovedWordsState struct{}

func (a ApprovedWordsState) Id() int {
	return APPROVED_WORDS_STATE_ID
}

func (a ApprovedWordsState) Routes() map[int]int {
	m := make(map[int]int)
	return m
}
