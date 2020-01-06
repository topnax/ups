package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"ups/sp/server/model"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
	"ups/sp/server/protocol/responses"
)

const (
	ConsecutiveInvalidMessageCountCeiling = 5
)

type KrisKrosRouter struct {
	callbacks                   map[int]func(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response
	server                      *KrisKrosServer
	Handlers                    []def.MessageHandler
	states                      map[int]State
	UserStates                  map[int]State
	SocketToUserID              map[int]int
	UserIDToSocket              map[int]int
	socketToInvalidMessageCount map[int]int
	IgnoreTransitionStateChange bool
	socketCloser                def.SocketCloser
	RouteMutex                  *sync.Mutex
}

// creates a new instance of a router
func newKrisKrosRouter(server *KrisKrosServer) *KrisKrosRouter {
	router := KrisKrosRouter{server: server}
	router.callbacks = make(map[int]func(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response)
	router.registerRoutes()
	router.registerStates()
	router.UserStates = make(map[int]State)
	router.SocketToUserID = make(map[int]int)
	router.UserIDToSocket = make(map[int]int)
	router.socketToInvalidMessageCount = make(map[int]int)
	router.RouteMutex = &sync.Mutex{}
	return &router
}

// registers a callback to a message handler
func (router *KrisKrosRouter) register(handler def.MessageHandler, callback func(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response) {
	router.callbacks[handler.GetType()] = callback
	router.Handlers = append(router.Handlers, handler)
}

// returns a response indicating that a message could not be cast
func failedToCast(message def.MessageHandler) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed to cast message of type %d", message.GetType()), impl.FailedToCast)
}

// registers a router state
func (router *KrisKrosRouter) registerState(state State) {
	state.Routes()[messages.KeepAliveMessageType] = state.Id()
	router.states[state.Id()] = state
}

// signals that an invalid message has been received, when N consecutive messages are received, connection to socket is closed
func (router *KrisKrosRouter) invalidMessage(socket int) {
	invalidMessages, exists := router.socketToInvalidMessageCount[socket]
	if exists {
		invalidMessages++
	} else {
		invalidMessages = 0
	}

	if invalidMessages >= ConsecutiveInvalidMessageCountCeiling {
		log.Warnf("Received %d invalid consecutive invalid message from SOCKET=%d, about to close connection to it", invalidMessages, socket)
		if router.socketCloser != nil {
			router.socketCloser.CloseFd(socket)
		} else {
			log.Errorf("Could not close SOCKET=%d, because socket closer is null", socket)
		}
		router.socketToInvalidMessageCount[socket] = 0
	} else {
		log.Warnf("Received an invalid message from SOCKET=%d, consecutive invalid message count is %d", socket, invalidMessages)
		router.socketToInvalidMessageCount[socket] = invalidMessages
	}
}

// signals that an valid message has been received,
func (router *KrisKrosRouter) validMessage(socket int) {
	invalidMessageCount, exists := router.socketToInvalidMessageCount[socket]
	if exists && invalidMessageCount > 0 {
		log.Infof("Resetting consecutive invalid message count of SOCKET=%d to 0", socket)
	}
	router.socketToInvalidMessageCount[socket] = 0
}

// routes a given message
func (router *KrisKrosRouter) route(message def.MessageHandler, socket int) def.Response {
	userID, exists := router.SocketToUserID[socket]

	var userState State

	// find an user based on it's socket
	if exists {
		userState = router.UserStates[userID]
	} else {
		userState = InitialState{}
		userID = -1
	}

	if message.GetType() == messages.UserLeavingMessageType {
		userState = InitialState{}
	}

	if message.GetType() == messages.KeepAliveMessageType {
		return keepAliveRoute(message, router.server, model.User{
			ID: socket,
		})
	}

	newStateID, exists := userState.Routes()[message.GetType()]

	if exists {
		log.Infof("Routing message of type %d when current state of user of id %d is of type %d", message.GetType(), userID, newStateID)
		route, exists := router.callbacks[message.GetType()]

		// find a route based on the message type
		if exists {
			// route if a route is found
			log.Infof("Route of type %d found", message.GetType())
			user, exists := router.server.usersById[userID]
			if !exists {
				user = &model.User{
					ID:   socket,
					Name: "NOT_CREATED",
				}
			}
			router.IgnoreTransitionStateChange = false
			response := route(message, router.server, *user)
			log.Infoln("responding:", response.Content())
			// check whether the route ended with a valid status
			if response.Type() < responses.ValidResponseCeiling || response.Type() == impl.PlainSuccess {
				if userID != -1 {
					state, exists := router.states[newStateID]
					if exists {
						// change user's state based on the route
						if router.IgnoreTransitionStateChange == false {
							router.UserStates[userID] = state
						} else {
							// route may change the state manually, don't overwrite
							log.Info("Ignored state transition...")
						}
						log.Infof("Routed message of type %d and switched to type %d", message.GetType(), state.Id())
					} else {
						log.Errorf("Could not get state from state map of ID %d", newStateID)
					}
				}
				router.validMessage(socket)
			} else {
				router.invalidMessage(socket)
			}
			return response
		}
		router.invalidMessage(socket)
		return impl.ErrorResponse(fmt.Sprintf("Could not route a message of type %d", message.GetType()), impl.FailedToRoute)
	}
	router.invalidMessage(socket)
	return impl.ErrorResponse(fmt.Sprintf("Cannot perform operation of type %d because current state of id %d does not allow it.", message.GetType(), userState.Id()), impl.OperationCannotBePerformed)
}

func (router *KrisKrosRouter) SetSocketCloser(socketCloser def.SocketCloser) {
	router.socketCloser = socketCloser
}
