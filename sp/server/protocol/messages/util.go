package messages

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
)

func parse(message def.Message, messageTemplate interface{}) bool {
	log.Infoln(message.Content())
	r := strings.NewReader(message.Content())
	d := json.NewDecoder(r)
	//err := json.Unmarshal([]byte(message.Content()), &messageTemplate)
	err := d.Decode(&messageTemplate)
	if err != nil {
		log.Errorf("JSON Unmarshal error: '%s'\nFrom message (type %d) of client #%d: '%s'", err, message.ClientID(), message.Type(), message.Content())
		log.Errorf("JSON Unmarshal error of content: '%s'", message.Content())
		return false
	}
	return true
}

func failedToParse(message def.Message) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed message of type %d, of content: '%s'", message.Type(), message.Content()), impl.FailedToParse)
}
