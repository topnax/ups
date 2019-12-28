package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"ups/sp/server/game"
	"ups/sp/server/model"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
	"ups/sp/server/protocol/responses"
)

type KrisKrosServer struct {
	count  int
	Router *KrisKrosRouter
	sender def.ResponseSender

	lobbyUIQ int
	userUIQ  int

	usersById   map[int]*model.User
	usersByName map[string]*model.User

	lobbies           map[int]*model.Lobby
	lobbiesByOwnerID  map[int]*model.Lobby
	lobbiesByPlayerID map[int]*model.Lobby

	gameServer *GameServer
}

func failedToCast(message def.MessageHandler) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed to cast message of type %d", message.GetType()), impl.FailedToCast)
}

func NewKrisKrosServer(sender def.ResponseSender) KrisKrosServer {
	kks := KrisKrosServer{
		lobbies:           make(map[int]*model.Lobby),
		lobbiesByOwnerID:  make(map[int]*model.Lobby),
		lobbiesByPlayerID: make(map[int]*model.Lobby),
		usersById:         make(map[int]*model.User),
		usersByName:       make(map[string]*model.User),
	}
	kks.gameServer = NewGameServer(&kks)
	kks.Router = newKrisKrosRouter(&kks)
	kks.sender = sender
	return kks
}

func (k KrisKrosServer) Read(message def.MessageHandler, clientUID int) def.Response {
	res := k.Router.route(message, clientUID)
	if res != nil {
		log.Infoln("Server returning: ", res.Content())
	} else {
		log.Infoln("Response from router is NIL")
	}
	return res
}

func (k KrisKrosServer) Send(response def.Response, userId int, responseId int) {
	socket, exists := k.Router.UserIDToSocket[userId]
	if exists {
		k.sender.Send(response, socket, responseId)
	} else {
		log.Errorf("Could not send to USERID=%d response of type %d and content '%s'", userId, response.Type(), response.Content())
	}
}

func (k KrisKrosServer) SendToPlayersOfState(response def.Response, targetStateID int, responseId int, userToBeIgnoredID int) {
	_, exists := k.Router.states[targetStateID]

	if exists {
		for userID, state := range k.Router.UserStates {
			if state.Id() == targetStateID && userID != userToBeIgnoredID {
				log.Debugf("SendToPlayersOfState targetStateID=%d sending to player of ID=%d", targetStateID, userID)
				k.Send(response, userID, responseId)
			}
		}
	} else {
		log.Errorf("Cannot send response to players of state of id %d, because such state does not exist.", targetStateID)
	}
}

func (k *KrisKrosServer) OnPlayerJoined(user model.User) def.Response {
	log.Infoln("On player joined:)", user.Name, k.count)
	k.count++
	return impl.SuccessResponse("Successfully :)")
}

func (k *KrisKrosServer) OnJoinLobby(message messages.JoinLobbyMessage, user model.User) def.Response {
	lobby, exists := k.lobbies[message.LobbyID]
	if exists {

		if len(lobby.Players) > 3 {
			return impl.ErrorResponse("Cannot join lobby. Player limit exceeded.", impl.LobbyPlayerLimitExceeded)
		}

		newPlayer := game.Player{
			Name:  user.Name,
			ID:    user.ID,
			Ready: false,
		}

		lobby.Players = append(lobby.Players, newPlayer)

		k.lobbiesByPlayerID[newPlayer.ID] = lobby

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
			if player.ID != newPlayer.ID {
				//resp := responses.PlayerJoinedResponse{PlayerName: newPlayer.Name, PlayerID: newPlayer.ID}
				resp := responses.LobbyUpdatedResponse{Lobby: *lobby}
				k.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
			}
		}

		// send notification to players looking for lobbies
		notif := k.OnGetLobbies(messages.GetLobbiesMessage{}, model.User{})
		k.SendToPlayersOfState(notif, AUTHORIZED_STATE_ID, 0, user.ID)

		resp := responses.LobbyJoinedResponse{Player: newPlayer, Lobby: *lobby}
		return impl.MessageResponse(resp, resp.Type())
	}
	return impl.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID), impl.LobbyDoesNotExist)
}

func (k *KrisKrosServer) OnCreateLobby(msg messages.CreateLobbyMessage, user model.User) def.Response {
	log.Infof("Receiver create message from user of ID %d", user.ID)
	_, exists := k.lobbiesByOwnerID[user.ID]
	if !exists {
		player := game.Player{
			Name:  user.Name,
			ID:    user.ID,
			Ready: false,
		}
		k.lobbies[k.lobbyUIQ] = &model.Lobby{
			Owner: player,
			ID:    k.lobbyUIQ,
		}
		lobby := k.lobbies[k.lobbyUIQ]
		k.lobbiesByPlayerID[player.ID] = lobby
		k.lobbiesByOwnerID[user.ID] = k.lobbies[k.lobbyUIQ]

		k.lobbyUIQ++

		lobby.Players = append(lobby.Players, player)

		// send notification to players looking for lobbies
		notif := k.OnGetLobbies(messages.GetLobbiesMessage{}, model.User{})
		k.SendToPlayersOfState(notif, AUTHORIZED_STATE_ID, 0, user.ID)

		resp := responses.LobbyJoinedResponse{Player: player, Lobby: *lobby}
		return impl.MessageResponse(resp, resp.Type())
	} else {
		return impl.ErrorResponse(fmt.Sprintf("Player of ID %d already created a lobby", user.ID), impl.PlayerAlreadyCreatedLobby)
	}
}

func (k *KrisKrosServer) OnGetLobbies(message messages.GetLobbiesMessage, user model.User) def.Response {
	lobbies := []model.Lobby{}

	for _, lobby := range k.lobbies {
		lobbies = append(lobbies, *lobby)
	}

	resp := responses.GetLobbiesResponse{Lobbies: lobbies}

	return impl.MessageResponse(resp, resp.Type())
}

func (k *KrisKrosServer) ClientDisconnected(socket int) {
	k.OnClientDisconnected(socket)
}

func (k *KrisKrosServer) removeClientFromLobby(userID int) bool {
	lobby, ok := k.lobbiesByPlayerID[userID]
	log.Infof("Inside remove client from lobby, user ID %d, lobby exists : %v", userID, ok)
	if ok {
		var dcdPlayer game.Player
		dcdPlayerIndex := -1
		for i, player := range lobby.Players {
			if player.ID == userID {
				dcdPlayer = player
				dcdPlayer.Ready = false
				dcdPlayerIndex = i
				break
			}
		}
		if dcdPlayer.ID == lobby.Owner.ID {
			k.destroyLobby(lobby)
		} else {

			lobby.Players = append(lobby.Players[:dcdPlayerIndex], lobby.Players[dcdPlayerIndex+1:]...)
			delete(k.lobbiesByPlayerID, userID)

			for _, player := range lobby.Players {
				if player.ID != userID {
					resp := responses.LobbyUpdatedResponse{Lobby: *lobby}
					k.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
				}
			}
		}

		// send notification to players looking for lobbies
		notif := k.OnGetLobbies(messages.GetLobbiesMessage{}, model.User{})
		k.SendToPlayersOfState(notif, AUTHORIZED_STATE_ID, 0, userID)

		return true
	}
	return false
}

func (k *KrisKrosServer) destroyLobby(lobby *model.Lobby) {
	log.Infof("Destroying a lobby of id %d and owner %s", lobby.ID, lobby.Owner.Name)
	for _, player := range lobby.Players {
		if player.ID != lobby.Owner.ID {
			resp := responses.LobbyDestroyedResponse{}
			delete(k.lobbiesByPlayerID, player.ID)
			k.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
			k.Router.UserStates[player.ID] = AuthorizedState{}
		}
	}

	delete(k.lobbiesByOwnerID, lobby.Owner.ID)
	delete(k.lobbies, lobby.ID)
}

func (k *KrisKrosServer) OnLeaveLobby(userID int) def.Response {
	if k.removeClientFromLobby(userID) {
		return impl.SuccessResponse("Successfully left lobby")
	} else {
		return impl.ErrorResponse("Could not leave the lobby", impl.CouldNotLeaveLobby)
	}
}

func (k *KrisKrosServer) OnStartLobby(userID int) def.Response {
	lobby, exists := k.lobbiesByOwnerID[userID]
	if exists && lobby.IsStartable() {
		resp := responses.LobbyStartedResponse{}
		out := impl.MessageResponse(resp, resp.Type())
		for _, player := range lobby.Players {
			delete(k.lobbiesByPlayerID, player.ID)
			if player.ID != lobby.Owner.ID {
				k.Send(out, player.ID, 0)
				k.Router.UserStates[player.ID] = GameStartedState{}
			}
		}

		delete(k.lobbiesByOwnerID, lobby.Owner.ID)
		delete(k.lobbies, lobby.ID)

		lobbies := []model.Lobby{}

		for _, lobby := range k.lobbies {
			lobbies = append(lobbies, *lobby)
		}

		k.SendToPlayersOfState(impl.StructMessageResponse(responses.GetLobbiesResponse{Lobbies: lobbies}), AUTHORIZED_STATE_ID, 0, -1)

		k.gameServer.CreateGame(lobby.Players)

		return out
	}
	return impl.ErrorResponse(fmt.Sprintf("Cannot create a lobby because such user of ID=%d is not an owner of a lobby", userID), impl.GeneralError)
}

func (k *KrisKrosServer) OnPlayerReadyToggle(playerID int, ready bool) def.Response {
	log.Infof("Setting %d to %v", playerID, ready)
	lobby, exists := k.lobbiesByPlayerID[playerID]
	if exists {
		found := false
		readyPlayerIndex := 0
		for index, player := range lobby.Players {
			if player.ID == playerID {
				found = true
				readyPlayerIndex = index
				break
			}
		}
		if found {
			lobby.Players[readyPlayerIndex].Ready = ready
			for _, player := range lobby.Players {
				if player.ID != playerID {
					resp := responses.LobbyUpdatedResponse{Lobby: *lobby}
					k.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
				}
			}
			resp := responses.LobbyUpdatedResponse{Lobby: *lobby}
			log.Infoln("Owner is", lobby.Owner.Ready)
			return impl.MessageResponse(resp, resp.Type())
		}
	}
	return impl.ErrorResponse("Could not find such user in a lobby", impl.CouldNotFindSuchUserInLobby)
}

func (k *KrisKrosServer) OnUserAuthenticate(clientUID int, message messages.UserAuthenticationMessage) def.Response {
	if len(strings.Trim(message.Name, " ")) == 0 {
		return impl.ErrorResponse(fmt.Sprintf("Name must not be empty"), impl.NameMustNotBeEmpty)
	}
	user, exists := k.usersByName[message.Name]

	if !exists {
		user = &model.User{
			Name: message.Name,
			ID:   k.userUIQ,
		}
		k.userUIQ++
		k.usersById[user.ID] = user
		k.usersByName[user.Name] = user
		k.Router.SocketToUserID[clientUID] = user.ID
		k.Router.UserIDToSocket[user.ID] = clientUID
		k.Router.UserStates[user.ID] = AuthorizedState{}
		resp := responses.UserAuthenticatedResponse{User: *user}
		return impl.MessageResponse(resp, resp.Type())
	} else {
		_, exists = k.Router.UserIDToSocket[user.ID]
		if !exists {
			log.Infof("An user of ID %d and name %s has reconnected via socket %d", user.ID, user.Name, clientUID)
			k.Router.SocketToUserID[clientUID] = user.ID
			k.Router.UserIDToSocket[user.ID] = clientUID
			resp := responses.UserAuthenticatedResponse{User: *user}
			return impl.MessageResponse(resp, resp.Type())
		}
	}

	return impl.ErrorResponse(fmt.Sprintf("User of name %s is already logged on the server under ID of %d", user.Name, user.ID), impl.PlayerNameAlreadyTaken)
}

func (k *KrisKrosServer) OnClientDisconnected(clientUID int) {
	userID, exists := k.Router.SocketToUserID[clientUID]
	log.Debugf("User of Socket ID %d has disconnected\n", clientUID)
	if exists {
		log.Infof("Deleting a socket %d and %d from UserIDToSocket map", clientUID, userID)
		delete(k.Router.SocketToUserID, clientUID)
		delete(k.Router.UserIDToSocket, userID)
		user, exists := k.usersById[userID]
		if exists {
			state, exists := k.Router.UserStates[user.ID]
			if exists {
				if state.Id() >= LOBBY_JOINED_ID && state.Id() <= LOBBY_CREATED_READY_ID {
					k.removeClientFromLobby(user.ID)
					delete(k.usersById, userID)
					delete(k.usersByName, user.Name)
					delete(k.Router.UserStates, user.ID)
					log.Infof("Deleting a player of name %s", user.Name)
				} else if state.Id() >= GAME_STARTED_STATE_ID && state.Id() <= WORDS_VALIDITY_DECIDED_STATE_ID {
					k.gameServer.PlayerDisconnected(userID, state.Id())
				} else {
					// TODO
				}
			}
		}
	}
}

func (k *KrisKrosServer) OnUserDisconnecting(clientUID int) {
	userID, exists := k.Router.SocketToUserID[clientUID]
	log.Debugf("User of Socket ID %d wants to leave\n", clientUID)
	if exists {
		user, exists := k.usersById[userID]
		if exists {
			delete(k.usersById, userID)
			delete(k.usersByName, user.Name)
			log.Infof("Deleting a player of name %s", user.Name)
		}
	}
}

//func (k *KrisKrosServer) OnGetLobbies(mes encoding.Message, clientUID int) encoding.ResponseMessage {
//	var lobbies []model.Lobby
//	for _, v := range k.lobbies {
//		lobbies = append(lobbies, *v)
//	}
//	k.SendMessage(LobbiesListMessage{Lobbies: lobbies}, clientUID)
//
//	return encoding.SuccessResponse("great")
//}
//
//type LobbiesListMessage struct {
//	Lobbies []model.Lobby `json:"lobbies"`
//}
//
//func (p LobbiesListMessage) GetType() int {
//	return 103
//}
