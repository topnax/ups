package impl

import (
	"fmt"
	"ups/sp/server/encoding"
	"ups/sp/server/messages"
)

type SimpleJsonReader struct {
	handlers                 map[int]encoding.Message
	applicationMessageReader encoding.ApplicationMessageReader
}

func (s *SimpleJsonReader) Init() {
	s.handlers = make(map[int]encoding.Message)
	s.Register(&messages.GetLobbiesMessage{})
	s.Register(&messages.CreateLobbyMessage{})
	s.Register(&messages.JoinLobbyMessage{})
}

func (s *SimpleJsonReader) Register(handler encoding.Message) {
	s.handlers[(handler).GetType()] = handler
}

func (s *SimpleJsonReader) Read(message encoding.SimpleMessage) encoding.ResponseMessage {
	handler, ok := s.handlers[message.Type]
	if !ok {
		return encoding.ErrorResponse(fmt.Sprintf("Cannot read message from client %d of type %d\nContent: '%s'", message.ClientUID, message.Type, message.Content))
	} else {
		return handler.Handle(message, s.applicationMessageReader)
	}
}

func (s *SimpleJsonReader) SetOutput(reader encoding.ApplicationMessageReader) {
	s.applicationMessageReader = reader
}
