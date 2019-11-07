package encoding

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type JsonReader interface {
	Read(message SimpleMessage) ResponseMessage
	SetOutput(reader ApplicationMessageReader)
}

type SimpleJsonReader struct {
	handlers                 map[int]Message
	applicationMessageReader ApplicationMessageReader
}

func (s *SimpleJsonReader) Init() {
	s.handlers = make(map[int]Message)
	s.Register(&CreateLobbyMessage{})
	s.Register(&JoinLobbyMessage{})
}

func (s *SimpleJsonReader) Register(handler Message) {
	s.handlers[(handler).GetType()] = handler
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

func (s *SimpleJsonReader) Read(message SimpleMessage) ResponseMessage {
	handler, ok := s.handlers[message.Type]
	if !ok {
		return ErrorResponse(fmt.Sprintf("Cannot read message from client %d of type %d\nContent: '%s'", message.ClientUID, message.Type, message.Content))
	} else {
		return handler.Handle(message, s.applicationMessageReader)
	}
}

func (s *SimpleJsonReader) SetOutput(reader ApplicationMessageReader) {
	s.applicationMessageReader = reader
}
