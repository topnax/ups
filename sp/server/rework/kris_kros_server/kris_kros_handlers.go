package kris_kros_server

import (
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/messages"
)

func (k *KrisKrosServer) GetHandlers() []def.MessageHandler {
	return []def.MessageHandler{
		&messages.PlayerJoinedMessage{},
	}
}
