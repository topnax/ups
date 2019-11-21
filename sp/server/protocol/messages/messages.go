package messages

import (
	"ups/sp/server/protocol/def"
)

const (
	PlayerJoinedLobbyType = 1
	CreateLobbyType       = 2
	GetLobbiesType        = 3
	JoinLobbyMessageType  = 4
	LeaveLobbyMessageType = 5
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
	PlayerName string `json:"player_name"`
}

func (p CreateLobbyMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		return amr.Read(p, message.ClientID())
	}
	return failedToParse(message)
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
	PlayerName string `json:"player_name"`
	LobbyID    int    `json:"lobby_id"`
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