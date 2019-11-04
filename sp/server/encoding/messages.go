package encoding

import "github.com/sirupsen/logrus"

func (s *SimpleJsonReader) Init() {
	s.handlers = GetMessageHandlers()
}

type CreatedMessageHandler struct {
	Surname string        `json:"surname"`
	Smr     SampleMessage `json:"smr"`
}

type SampleMessage struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c *CreatedMessageHandler) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&c) {
		logrus.Debugf("CreatedMessageReceived, surname %s,name %s, age %d", c.Surname, c.Smr.Name, c.Smr.Age)
	}
}

func (s *SampleMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&s) {
		logrus.Debugf("Simple message received, name %s, age %d", s.Name, s.Age)
	}
}

type CreateLobbyMessage struct {
	ClientID int
}

func (c *CreateLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&c) {
		amr.OnCreateLobby(*c, message.ClientUID)
	}
}

type JoinLobbyMessage struct {
	LobbyID int `json:"lobby_id"`
}

func (c *JoinLobbyMessage) Handle(message SimpleMessage, amr ApplicationMessageReader) {
	if message.Parse(&c) {
		amr.OnJoinLobby(*c, message.ClientUID)
	}
}
