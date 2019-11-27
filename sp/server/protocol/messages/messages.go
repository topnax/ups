package messages

import (
	"ups/sp/server/protocol/def"
)

const (
	PlayerJoinedLobbyType         = 1
	CreateLobbyType               = 2
	GetLobbiesType                = 3
	JoinLobbyMessageType          = 4
	LeaveLobbyMessageType         = 5
	PlayerReadyMessageType        = 6
	UserAuthenticationMessageType = 7
	UserLeavingMessageType        = 8
)

type PlayerJoinedMessage struct {
	PlayerName string `json:"player_name"`
}

func (p PlayerJoinedMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
}

func (p PlayerJoinedMessage) GetType() int {
	return PlayerJoinedLobbyType
}

////////////////////////////////////////////

type CreateLobbyMessage struct {
}

func (p CreateLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	return amr.Read(p, message.ClientID())
}

func (p CreateLobbyMessage) GetType() int {
	return CreateLobbyType
}

////////////////////////////////////////////

type GetLobbiesMessage struct {
}

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
	Name string `json:"name"`
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
