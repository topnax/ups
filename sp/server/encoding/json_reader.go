package encoding

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type MessageHandler interface {
	Handle(message SimpleMessage, amr ApplicationMessageReader)
	GetType() int
}

type JsonReader interface {
	Read(message SimpleMessage)
	SetOutput(reader ApplicationMessageReader)
}

type SimpleJsonReader struct {
	handlers                 map[int]MessageHandler
	applicationMessageReader ApplicationMessageReader
}

func (s *SimpleJsonReader) Init() {
	s.handlers = make(map[int]MessageHandler)
	s.Register(&CreateLobbyMessage{})
	s.Register(&JoinLobbyMessage{})
}

func (s *SimpleJsonReader) Register(handler MessageHandler) {
	s.handlers[(handler).GetType()] = handler
}

func (simpleMessage SimpleMessage) Parse(messageTemplate interface{}) bool {
	logrus.Infoln(simpleMessage.Content)
	err := json.Unmarshal([]byte(simpleMessage.Content), &messageTemplate)
	if err != nil {
		logrus.Errorf("JSON Unmarshal error: '%s'\nFrom message of client #%d: '%s'", err, simpleMessage.ClientUID, simpleMessage.Content)
		return false
	}
	return true
}

func (s *SimpleJsonReader) Read(message SimpleMessage) {
	handler, ok := s.handlers[message.Type]
	if !ok {
		logrus.Errorf("Cannot read message from client %d of type %d\nContent: '%s'", message.ClientUID, message.Type, message.Content)
	} else {
		handler.Handle(message, s.applicationMessageReader)
	}
}

func (s *SimpleJsonReader) SetOutput(reader ApplicationMessageReader) {
	s.applicationMessageReader = reader
}
