package encoding

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type JsonReader interface {
	Read(message SimpleMessage) ResponseMessage
	SetOutput(reader ApplicationMessageReader)
}

type Message interface {
	Handle(message SimpleMessage, amr ApplicationMessageReader) ResponseMessage
	GetType() int
}

type TypedMessage interface {
	GetType() int
}

func (simpleMessage SimpleMessage) Parse(messageTemplate interface{}) bool {
	logrus.Infoln(simpleMessage.Content)
	err := json.Unmarshal([]byte(simpleMessage.Content), &messageTemplate)
	if err != nil {
		logrus.Errorf("JSON Unmarshal error: '%s'\nFrom message (type %d) of client #%d: '%s'", err, simpleMessage.ClientUID, simpleMessage.Type, simpleMessage.Content)
		return false
	}
	return true
}
