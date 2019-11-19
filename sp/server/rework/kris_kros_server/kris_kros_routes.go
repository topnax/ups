package kris_kros_server

import (
	log "github.com/sirupsen/logrus"
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/messages"
)

func (router *KrisKrosRouter) registerRoutes() {
	router.register(messages.PlayerJoinedMessage{}, playerJoinedRoute)
}

func playerJoinedRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	// $15#1#1#{"name":"Stan"}
	log.Infoln("Player joined route triggered")

	msg, ok := handler.(messages.PlayerJoinedMessage)

	if ok {
		return server.OnPlayerJoined(msg, clientUID)
	}
	return failedToCast(handler)
}
