package kris_kros_server

import (
	"time"
	"ups/sp/server/model"
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
	router.register(messages.StartLobbyMessage{}, startLobbyRoute)
	router.register(messages.LetterPlacedMessage{}, letterPlacedRoute)
	router.register(messages.LetterRemovedMessage{}, letterRemovedRoute)
	router.register(messages.FinishRoundMessage{}, finishRoundRoute)
	router.register(messages.ApproveWordsMessage{}, approveWordsRoute)
	router.register(messages.DeclineWordsMessage{}, declineWordsRoute)
	router.register(messages.KeepAliveMessage{}, keepAliveRoute)
	router.register(messages.LeaveGameMessage{}, leaveGameRoute)
}

func playerJoinedRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.OnPlayerJoined(user)
}

func createLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.CreateLobbyMessage)

	if ok {
		return server.OnCreateLobby(msg, user)
	}

	return failedToCast(handler)
}

func getLobbiesRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.GetLobbiesMessage)

	if ok {
		return server.OnGetLobbies(msg, user)
	}

	return failedToCast(handler)
}

func joinLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.JoinLobbyMessage)

	if ok {
		return server.OnJoinLobby(msg, user)
	}

	return failedToCast(handler)
}

func leaveLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	_, ok := handler.(messages.LeaveLobbyMessage)

	if ok {
		return server.OnLeaveLobby(user.ID)
	}

	return failedToCast(handler)
}

func playerReadyRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.PlayerReadyToggle)

	if ok {
		return server.OnPlayerReadyToggle(user.ID, msg.Ready)
	}

	return failedToCast(handler)
}

func userAuthenticationRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.UserAuthenticationMessage)

	if ok {
		return server.OnUserAuthenticate(user.ID, msg)
	}

	return failedToCast(handler)
}

func userLeavingRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	server.OnUserDisconnecting(user.ID)
	return impl.DoNotRespond()
}

func startLobbyRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.OnStartLobby(user.ID)
}

func letterPlacedRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.LetterPlacedMessage)

	if ok {
		return server.gameServer.OnLetterPlaced(user.ID, msg)
	}

	return failedToCast(handler)
}

func letterRemovedRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	msg, ok := handler.(messages.LetterRemovedMessage)

	if ok {
		return server.gameServer.OnLetterRemoved(user.ID, msg)
	}

	return failedToCast(handler)
}

func finishRoundRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.gameServer.OnFinishRound(user.ID)
}

func approveWordsRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.gameServer.OnApproveWords(user.ID)
}

func declineWordsRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.gameServer.OnDeclineWords(user.ID)
}

func keepAliveRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	server.userLastKeepAlive[user.ID] = time.Now()
	return server.onKeepAlive(user.ID)
}

func leaveGameRoute(handler def.MessageHandler, server *KrisKrosServer, user model.User) def.Response {
	return server.gameServer.OnPlayerLeavingGame(user.ID)
}
