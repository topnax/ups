package kris_kros_server

import (
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/messages"
)

func (k *KrisKrosServer) GetHandlers() []def.MessageHandler {
	return []def.MessageHandler{
		&messages.PlayerJoinedMessage{},
	}
}
