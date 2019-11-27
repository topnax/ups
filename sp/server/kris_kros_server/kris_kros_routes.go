package kris_kros_server

import (
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
)

func (router *KrisKrosRouter) registerRoutes() {
	router.register(messages.PlayerJoinedMessage{}, playerJoinedRoute)
	router.register(messages.CreateLobbyMessage{}, createLobbyRoute)
	router.register(messages.GetLobbiesMessage{}, getLobbiesRoute)
	router.register(messages.JoinLobbyMessage{}, joinLobbyRoute)
	router.register(messages.LeaveLobbyMessage{}, leaveLobbyRoute)
	router.register(messages.PlayerReadyToggle{}, playerReadyRoute)
	router.register(messages.UserAuthenticationMessage{}, userAuthenticationRoute)
	router.register(messages.UserLeavingMessage{}, userLeavingRoute)
}

func playerJoinedRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.PlayerJoinedMessage)

	if ok {
		return server.OnPlayerJoined(msg, clientUID)
	}

	return failedToCast(handler)
}

func createLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.CreateLobbyMessage)

	if ok {
		return server.OnCreateLobby(msg, clientUID)
	}

	return failedToCast(handler)
}

func getLobbiesRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.GetLobbiesMessage)

	if ok {
		return server.OnGetLobbies(msg, clientUID)
	}

	return failedToCast(handler)
}

func joinLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.JoinLobbyMessage)

	if ok {
		return server.OnJoinLobby(msg, clientUID)
	}

	return failedToCast(handler)
}

func leaveLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	_, ok := handler.(messages.LeaveLobbyMessage)

	if ok {
		return server.OnLeaveLobby(clientUID)
	}

	return failedToCast(handler)
}

func playerReadyRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.PlayerReadyToggle)

	if ok {
		return server.OnPlayerReadyToggle(clientUID, msg.Ready)
	}

	return failedToCast(handler)
}

func userAuthenticationRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	msg, ok := handler.(messages.UserAuthenticationMessage)

	if ok {
		return server.OnUserAuthenticate(clientUID, msg)
	}

	return failedToCast(handler)
}

func userLeavingRoute(handler def.MessageHandler, server *KrisKrosServer, clientUID int) def.Response {
	server.OnUserDisconnecting(clientUID)
	return impl.DoNotRespond()
}
