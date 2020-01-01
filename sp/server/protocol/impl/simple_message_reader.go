package impl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/protocol/def"
)

type SimpleMessageReader struct {
	handlers                 map[int]def.MessageHandler
	applicationMessageReader def.ApplicationMessageReader
}

// creates a new message reader based on the slice of handlers passed in
func NewSimpleMessageReader(reader def.ApplicationMessageReader, handlers []def.MessageHandler) SimpleMessageReader {
	smr := SimpleMessageReader{}
	smr.handlers = make(map[int]def.MessageHandler)
	smr.applicationMessageReader = reader

	for _, handler := range handlers {
		smr.Register(handler)
	}

	return smr
}

// registers a message handler
func (s *SimpleMessageReader) Register(handler def.MessageHandler) {
	s.handlers[handler.GetType()] = handler
}

// reads a parsed message and uses a handler to handle it
func (s *SimpleMessageReader) Read(message def.Message) def.Response {
	handler, ok := s.handlers[message.Type()]
	log.Debugln("MessageReader read message from UID %d of type %d and of content %s", message.ClientID(), message.Type(), message.Content())
	if !ok {
		res := ErrorResponseID(fmt.Sprintf("Could not find a message handler for a message of type '%d' and content '%s'", message.Type(), message.Content()), NoMessageHandler, message.ID())
		return res
	} else {
		res := handler.Handle(message, s.applicationMessageReader)
		return res
	}
}
