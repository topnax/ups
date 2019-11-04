package encoding

type CreateLobbyMessage struct {
	ClientName string `json:"client_name"`
}

func (c *CreateLobbyMessage) GetType() int {
	return 1
}

func (c *CreateLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&c) {
		amr.OnCreateLobby(*c, message.ClientUID)
	}
}

func (c *JoinLobbyMessage) GetType() int {
	return 2
}

type JoinLobbyMessage struct {
	LobbyID    int    `json:"lobby_id"`
	ClientName string `json:"client_name"`
}

func (c *JoinLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&c) {
		amr.OnJoinLobby(*c, message.ClientUID)
	}
}
