package messages

import (
	"github.com/sirupsen/logrus"
	"ups/sp/server/model"
	"ups/sp/server/protocol/def"
)

const (
	PlayerJoinedLobby = 1
)

type PlayerJoinedMessage struct {
	Name string `json:"name"`
}

func (p PlayerJoinedMessage) Handle(message def.Message, amr def.ApplicationMessageReader) def.Response {
	if parse(message, &p) {
		logrus.Infoln("My name is", p.Name)
		res := amr.Read(p, message.ClientID())
		logrus.Infoln("MEssages returning", res.Content())
		return res
	}
	return failedToParse(message)
}

func (p PlayerJoinedMessage) GetType() int {
	return PlayerJoinedLobby
}

type PlayerConfirmation struct {
	counter     int
	model.Lobby `json:"lobby"`
}
