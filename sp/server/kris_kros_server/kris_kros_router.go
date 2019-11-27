package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/responses"
)

type KrisKrosRouter struct {
	callbacks      map[int]func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response
	server         *KrisKrosServer
	Handlers       []def.MessageHandler
	states         map[int]State
	UserStates     map[int]State
	SocketToUserID map[int]int
	UserIDToSocket map[int]int
}

func (router *KrisKrosRouter) register(handler def.MessageHandler, callback func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response) {
	router.callbacks[handler.GetType()] = callback
	router.Handlers = append(router.Handlers, handler)
}

func (router *KrisKrosRouter) registerState(state State) {
	router.states[state.Id()] = state
}

func (router *KrisKrosRouter) route(message def.MessageHandler, clientUID int) def.Response {
	userID, exists := router.SocketToUserID[clientUID]

	var userState State

	if exists {
		userState = router.UserStates[userID]
	} else {
		userState = InitialState{}
		userID = -1
	}

	newStateID, exists := userState.Routes()[message.GetType()]

	if exists {
		log.Infof("Routing message of type %d when current state of user of id %d is of type %s", message.GetType(), userID, newStateID)
		route, exists := router.callbacks[message.GetType()]

		if exists {
			log.Infof("Route of type %d found", message.GetType())
			response := route(message, router.server, clientUID)
			log.Infoln("responding:", response.Content())
			if response.Type() < responses.ValidResponseCeiling {
				if userID != -1 {
					state, exists := router.states[newStateID]
					if exists {
						router.UserStates[userID] = state
						log.Infof("Routed message of type %d and switched to type %ď", message.GetType(), state)
					} else {
						log.Errorf("Could not get state from state map of ID %d", state)
					}
				}
			}
			return response
		}

		return impl.ErrorResponse(fmt.Sprintf("Could not route a message of type %d", message.GetType()), impl.FailedToRoute)
	}

	return impl.ErrorResponse(fmt.Sprintf("Cannot perform operation of type %d because current state of id %d does not allow it.", message.GetType(), userState.Id()), impl.OperationCannotBePerformed)
}

func newKrisKrosRouter(server *KrisKrosServer) KrisKrosRouter {
	router := KrisKrosRouter{server: server}
	router.callbacks = make(map[int]func(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response)
	router.registerRoutes()
	router.registerStates()
	router.UserStates = make(map[int]State)
	router.SocketToUserID = make(map[int]int)
	router.UserIDToSocket = make(map[int]int)
	return router
}
