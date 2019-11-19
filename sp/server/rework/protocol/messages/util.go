package messages

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/impl"
)

func parse(message def.Message, messageTemplate interface{}) bool {
	log.Infoln(message.Content())
	err := json.Unmarshal([]byte(message.Content()), &messageTemplate)
	if err != nil {
		log.Errorf("JSON Unmarshal error: '%s'\nFrom message (type %d) of client #%d: '%s'", err, message.ClientID(), message.Type(), message.Content())
		return false
	}
	return true
}

func failedToParse(message def.Message) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed message of type %d, of content: '%s'", message.Type, message.Content), impl.FailedToParse)
}
