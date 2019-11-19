package impl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/rework/protocol/def"
)

type SimpleMessageReader struct {
	handlers                 map[int]def.MessageHandler
	applicationMessageReader def.ApplicationMessageReader
}

func NewSimpleMessageReader(reader def.ApplicationMessageReader, handlers []def.MessageHandler) SimpleMessageReader {
	smr := SimpleMessageReader{}
	smr.handlers = make(map[int]def.MessageHandler)
	smr.applicationMessageReader = reader

	for _, handler := range handlers {
		smr.Register(handler)
	}

	return smr
}

func (s *SimpleMessageReader) Register(handler def.MessageHandler) {
	s.handlers[handler.GetType()] = handler
}

func (s *SimpleMessageReader) Read(message def.Message) def.Response {
	handler, ok := s.handlers[message.Type()]
	log.Debugln("MessageReader read message from UID %d of type %d and of content %s", message.ClientID(), message.Type(), message.Content())
	if !ok {
		return ErrorResponse(fmt.Sprintf("Could not find a message handler for a message of type '%d' and content '%s'", message.Type(), message.Content()), NoMessageHandler)
	} else {
		res := handler.Handle(message, s.applicationMessageReader)
		return res
	}
}
