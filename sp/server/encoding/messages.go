package encoding

import "fmt"

type Message interface {
	Handle(message SimpleMessage, amr ApplicationMessageReader) ResponseMessage
	GetType() int
}

type TypedMessage interface {
	GetType() int
}

func failedToParse(message SimpleMessage) ResponseMessage {
	return ErrorResponse(fmt.Sprintf("Failed message of type %d, of content: '%s'", message.Type, message.Content))
}

// created lobby

type CreateLobbyMessage struct {
	ClientName string `json:"client_name"`
}

func (c *CreateLobbyMessage) GetType() int {
	return 1
}

func (c *CreateLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) ResponseMessage {
	if message.Parse(&c) {
		return amr.OnCreateLobby(*c, message.ClientUID)
	}
	return failedToParse(message)
}

// join lobby

func (c *JoinLobbyMessage) GetType() int {
	return 2
}

type JoinLobbyMessage struct {
	LobbyID    int    `json:"lobby_id"`
	ClientName string `json:"client_name"`
}

func (c *JoinLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) ResponseMessage {
	if message.Parse(&c) {
		return amr.OnJoinLobby(*c, message.ClientUID)
	}
	return failedToParse(message)
}

// output

type PlayerJoinedLobbyMessage struct {
	ClientName string `json:"client_name"`
}

func (p PlayerJoinedLobbyMessage) GetType() int {
	return 101
}

type PlayerLeftLobbyMessage struct {
	ClientName string `json:"client_name"`
}

func (p PlayerLeftLobbyMessage) GetType() int {
	return 102
}
