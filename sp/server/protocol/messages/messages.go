package messages

import (
	"ups/sp/server/game"
	"ups/sp/server/protocol/def"
)

const (
	CreateLobbyType               = 2
	GetLobbiesType                = 3
	JoinLobbyMessageType          = 4
	LeaveLobbyMessageType         = 5
	PlayerReadyMessageType        = 6
	UserAuthenticationMessageType = 7
	UserLeavingMessageType        = 8
	StartLobbyMessageType         = 9
	LetterPlacedMessageType       = 10
	LetterRemovedMessageType      = 11
	FinishRoundMessageType        = 12
	ApproveWordsMessageType       = 13
	DeclineWordsMessageType       = 14
	KeepAliveMessageType          = 15
	LeaveGameMessageType          = 16
)

////////////////////////////////////////////

type CreateLobbyMessage struct{}

func (p CreateLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p CreateLobbyMessage) GetType() int {
	return CreateLobbyType
}

////////////////////////////////////////////

type GetLobbiesMessage struct{}

func (p GetLobbiesMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p GetLobbiesMessage) GetType() int {
	return GetLobbiesType
}

////////////////////////////////////////////

type JoinLobbyMessage struct {
	LobbyID int `json:"lobby_id"`
}

func (p JoinLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p JoinLobbyMessage) GetType() int {
	return JoinLobbyMessageType
}

////////////////////////////////////////////

type LeaveLobbyMessage struct{}

func (p LeaveLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p LeaveLobbyMessage) GetType() int {
	return LeaveLobbyMessageType
}

////////////////////////////////////////////

type PlayerReadyToggle struct {
	Ready bool `json:"ready"`
}

func (p PlayerReadyToggle) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p PlayerReadyToggle) GetType() int {
	return PlayerReadyMessageType
}

////////////////////////////////////////////

type UserAuthenticationMessage struct {
	Name         string `json:"name"`
	Reconnecting bool   `json:"reconnecting"`
}

func (p UserAuthenticationMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p UserAuthenticationMessage) GetType() int {
	return UserAuthenticationMessageType
}

////////////////////////////////////////////

type UserLeavingMessage struct{}

func (p UserLeavingMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p UserLeavingMessage) GetType() int {
	return UserLeavingMessageType
}

////////////////////////////////////////////

type StartLobbyMessage struct{}

func (p StartLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p StartLobbyMessage) GetType() int {
	return StartLobbyMessageType
}

////////////////////////////////////////////

type LetterPlacedMessage struct {
	Letter game.Letter `json:"letter"`
	Row    int         `json:"row"`
	Column int         `json:"column"`
}

func (p LetterPlacedMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p LetterPlacedMessage) GetType() int {
	return LetterPlacedMessageType
}

////////////////////////////////////////////

type LetterRemovedMessage struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

func (p LetterRemovedMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p LetterRemovedMessage) GetType() int {
	return LetterRemovedMessageType
}

////////////////////////////////////////////

type FinishRoundMessage struct{}

func (p FinishRoundMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p FinishRoundMessage) GetType() int {
	return FinishRoundMessageType
}

////////////////////////////////////////////

type ApproveWordsMessage struct{}

func (p ApproveWordsMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p ApproveWordsMessage) GetType() int {
	return ApproveWordsMessageType
}

////////////////////////////////////////////

type DeclineWordsMessage struct{}

func (p DeclineWordsMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p DeclineWordsMessage) GetType() int {
	return DeclineWordsMessageType
}

////////////////////////////////////////////

type KeepAliveMessage struct{}

func (p KeepAliveMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p KeepAliveMessage) GetType() int {
	return KeepAliveMessageType
}

////////////////////////////////////////////

type LeaveGameMessage struct{}

func (p LeaveGameMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p LeaveGameMessage) GetType() int {
	return LeaveGameMessageType
}
