package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/impl"
)

type KrisKrosRouter struct {
	callbacks map[int]func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response
	server    *KrisKrosServer
	Handlers  []def.MessageHandler
}

func (router *KrisKrosRouter) register(handler def.MessageHandler, callback func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response) {
	router.callbacks[handler.GetType()] = callback
	router.Handlers = append(router.Handlers, handler)
}

func (router *KrisKrosRouter) route(message def.MessageHandler, clientUID int) def.Response {
	log.Infof("Routing message of type %d", message.GetType())

	route, exists := router.callbacks[message.GetType()]

	if exists {
		log.Infof("Route of type %d found", message.GetType())
		response := route(message, router.server, clientUID)
		log.Infoln("responding:", response.Content())
		return response
	}

	return impl.ErrorResponse(fmt.Sprintf("Could not route a message of type %d", message.GetType()), impl.FailedToRoute)
}

func newKrisKrosRouter(server *KrisKrosServer) KrisKrosRouter {
	router := KrisKrosRouter{server: server}
	router.callbacks = make(map[int]func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response)
	router.registerRoutes()
	return router
}
