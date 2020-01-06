package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
	"ups/sp/server/game"
	"ups/sp/server/model"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
	"ups/sp/server/protocol/responses"
)

const (
	KeepAliveTimerPeriod = 3000
	MaximalKeepAliveAge  = 2500
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

	userLastKeepAlive      map[int]time.Time
	userLastKeepAliveMutex *sync.Mutex

	lobbyLimit int

	gameServer *GameServer
}

// creates a new kriskros server
func NewKrisKrosServer(sender def.ResponseSender, lobbyLimit int) KrisKrosServer {
	kks := KrisKrosServer{
		lobbies:                make(map[int]*model.Lobby),
		lobbiesByOwnerID:       make(map[int]*model.Lobby),
		lobbiesByPlayerID:      make(map[int]*model.Lobby),
		usersById:              make(map[int]*model.User),
		usersByName:            make(map[string]*model.User),
		userLastKeepAlive:      make(map[int]time.Time),
		userLastKeepAliveMutex: &sync.Mutex{},
		lobbyLimit:             lobbyLimit,
	}
	log.Infof("Current lobby limit is of %d", lobbyLimit)
	kks.gameServer = NewGameServer(&kks)
	kks.Router = newKrisKrosRouter(&kks)
	kks.sender = sender

	ticker := time.NewTicker(KeepAliveTimerPeriod * time.Millisecond)
	done := make(chan bool)

	// launch a goroutine that checks whether keep alive messages from clients are periodically coming
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				kks.userLastKeepAliveMutex.Lock()
				socketsThatDisconnected := []int{}
				for socket, lastKeepAlive := range kks.userLastKeepAlive {
					userID, exists := kks.Router.SocketToUserID[socket]
					if lastKeepAlive.Add(MaximalKeepAliveAge * time.Millisecond).Before(time.Now()) {
						socketsThatDisconnected = append(socketsThatDisconnected, socket)
						if exists {
							log.Warnf("User of ID=%d and SOCKET=%d disconnected via KEEP_ALIVE", userID, socket)
							kks.Router.RouteMutex.Lock()
							kks.OnClientDisconnected(socket, false)
							kks.Router.RouteMutex.Unlock()
						} else {
							log.Warnf("SOCKET=%d did not send KEEP_ALIVE in time but no user exists under such socket.", socket)
						}
					} else {
						log.Debugf("User of ID=%d and SOCKET=%d is OK", userID, socket)
					}
				}

				for _, socket := range socketsThatDisconnected {
					delete(kks.userLastKeepAlive, socket)
				}

				kks.userLastKeepAliveMutex.Unlock()
			}
		}
	}()

	return kks
}

// reads a message, passing it to the router to route it
func (k KrisKrosServer) Read(message def.MessageHandler, clientUID int) def.Response {
	k.Router.RouteMutex.Lock()
	res := k.Router.route(message, clientUID)
	k.Router.RouteMutex.Unlock()
	if res != nil {
		log.Infoln("Server returning: ", res.Content())
	} else {
		log.Infoln("Response from router is NIL")
	}
	return res
}

// sends a response to a user of the given ID
func (k KrisKrosServer) Send(response def.Response, userId int, responseId int) {
	socket, exists := k.Router.UserIDToSocket[userId]
	if exists {
		k.sender.Send(response, socket, responseId)
	} else {
		log.Warnf("Could not send to USERID=%d response of type %d and content '%s'", userId, response.Type(), response.Content())
	}
}

// sends a response to all users of the given state
func (k KrisKrosServer) SendToPlayersOfState(response def.Response, targetStateID int, responseId int, userToBeIgnoredID int) {
	_, exists := k.Router.states[targetStateID]

	if exists {
		for userID, state := range k.Router.UserStates {
			_, exists = k.Router.UserIDToSocket[userID]
			if exists {
				if state.Id() == targetStateID && userID != userToBeIgnoredID {
					log.Debugf("SendToPlayersOfState targetStateID=%d sending to player of ID=%d", targetStateID, userID)
					k.Send(response, userID, responseId)
				}
			}
		}
	} else {
		log.Errorf("Cannot send response to players of state of ID=%d, because such state does not exist.", targetStateID)
	}
}

// a handler of OnJoinLobby events
func (k *KrisKrosServer) OnJoinLobby(message messages.JoinLobbyMessage, user model.User) def.Response {
	lobby, exists := k.lobbies[message.LobbyID]
	if exists {
		log.Infof("OnJoinLobby by user of ID=%d", user.ID)

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
				resp := responses.LobbyUpdatedResponse{Lobby: *lobby}
				k.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
			}
		}

		// send notification to players looking for lobbies
		notification := k.OnGetLobbies(messages.GetLobbiesMessage{}, model.User{})
		k.SendToPlayersOfState(notification, AUTHORIZED_STATE_ID, 0, user.ID)

		resp := responses.LobbyJoinedResponse{Player: newPlayer, Lobby: *lobby}
		return impl.MessageResponse(resp, resp.Type())
	}
	return impl.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID), impl.LobbyDoesNotExist)
}

// a handler of OnCreateLobby events
func (k *KrisKrosServer) OnCreateLobby(msg messages.CreateLobbyMessage, user model.User) def.Response {
	log.Infof("OnCreateLobby by user of ID=%d", user.ID)
	_, exists := k.lobbiesByOwnerID[user.ID]
	if !exists {
		if len(k.lobbies) >= k.lobbyLimit {
			log.Warnln(fmt.Sprintf("Lobby limit exceeded, current lobby count is %d and the limit is %d", len(k.lobbies), k.lobbyLimit))
			return impl.ErrorResponse(fmt.Sprintf("Lobby limit exceeded, current lobby count is %d and the limit is %d", len(k.lobbies), k.lobbyLimit), impl.LobbyLimitExceeded)
		}
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

// handler of OnGetLobbies events
func (k *KrisKrosServer) OnGetLobbies(message messages.GetLobbiesMessage, user model.User) def.Response {
	lobbies := []model.Lobby{}

	for _, lobby := range k.lobbies {
		lobbies = append(lobbies, *lobby)
	}

	resp := responses.GetLobbiesResponse{Lobbies: lobbies}

	return impl.MessageResponse(resp, resp.Type())
}

// removes a user from a lobby
func (k *KrisKrosServer) removePlayerFromLobby(userID int) bool {
	lobby, ok := k.lobbiesByPlayerID[userID]
	log.Infof("Removing a player from lobby, user ID=%d, lobby exists : %v", userID, ok)
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
		notification := k.OnGetLobbies(messages.GetLobbiesMessage{}, model.User{})
		k.SendToPlayersOfState(notification, AUTHORIZED_STATE_ID, 0, userID)

		return true
	}
	return false
}

func (k *KrisKrosServer) destroyLobby(lobby *model.Lobby) {
	log.Infof("Destroying a lobby of ID=%d and owner NAME=%s", lobby.ID, lobby.Owner.Name)
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

// handler of OnLeaveLobby events
func (k *KrisKrosServer) OnLeaveLobby(userID int) def.Response {
	if k.removePlayerFromLobby(userID) {
		return impl.SuccessResponse("Successfully left lobby")
	} else {
		return impl.ErrorResponse("Could not leave the lobby", impl.CouldNotLeaveLobby)
	}
}

// handler of OnStartLobby events
func (k *KrisKrosServer) OnStartLobby(userID int) def.Response {
	lobby, exists := k.lobbiesByOwnerID[userID]
	if exists && lobby.IsStartPossible() {
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

		// send update to players looking for lobbies
		k.SendToPlayersOfState(impl.StructMessageResponse(responses.GetLobbiesResponse{Lobbies: lobbies}), AUTHORIZED_STATE_ID, 0, -1)

		k.gameServer.CreateGame(lobby.Players)

		return out
	}
	return impl.ErrorResponse(fmt.Sprintf("Cannot create a lobby because such user of ID=%d is not an owner of a lobby", userID), impl.GeneralError)
}

// handler of OnPlayerReadyToggle messages
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

// handler of OnUserAuthenticate events
func (k *KrisKrosServer) OnUserAuthenticate(clientUID int, message messages.UserAuthenticationMessage) def.Response {
	if len(strings.Trim(message.Name, " ")) == 0 {
		return impl.ErrorResponse(fmt.Sprintf("Name must not be empty"), impl.NameMustNotBeEmpty)
	}
	user, exists := k.usersByName[message.Name]

	// check whether user under such name already exists on the server
	if !exists {
		// such user does not exist, create a new one
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
		var resp responses.TypedResponse
		if message.Reconnecting {
			resp = responses.UserStateRegeneration{
				State: responses.ServerRestarted,
				User:  *user,
			}
		} else {
			resp = responses.UserAuthenticatedResponse{User: *user}
		}

		return impl.StructMessageResponse(resp)
	} else {
		// the socket is trying to reconnect - user with such name exists
		log.Warnf("User of socket=%d is trying go reconnect to userID=%d and username=%s", clientUID, user.ID, user.Name)
		_, exists = k.Router.UserIDToSocket[user.ID]
		if !exists {
			// no socket is connected to the name, we can reconnect the user
			log.Warnf("An user of ID %d and name %s has reconnected via socket %d", user.ID, user.Name, clientUID)
			k.Router.SocketToUserID[clientUID] = user.ID
			k.Router.UserIDToSocket[user.ID] = clientUID
			resp := responses.UserAuthenticatedResponse{User: *user}
			state, exists := k.Router.UserStates[user.ID]
			if exists {
				// state regeneration of the user
				if state.Id() >= GAME_STARTED_STATE_ID && state.Id() <= WORDS_VALIDITY_DECIDED_STATE_ID {
					regresp := k.gameServer.PlayerReconnected(user.ID)
					if regresp != nil {
						return regresp
					} else {
						return impl.StructMessageResponse(resp)
					}
				} else if message.Reconnecting {
					if state.Id() >= LOBBY_JOINED_ID && state.Id() <= LOBBY_CREATED_READY_ID {
						k.Router.UserStates[user.ID] = AuthorizedState{}
						return impl.StructMessageResponse(responses.UserStateRegeneration{
							State: responses.MovedToLobbyScreen,
							User:  *user,
						})
					} else if state.Id() != -1 {
						k.Router.UserStates[user.ID] = AuthorizedState{}
						return impl.StructMessageResponse(responses.UserStateRegeneration{
							State: responses.Nothing,
							User:  *user,
						})
					}
				}
			}
			k.Router.UserStates[user.ID] = AuthorizedState{}

			// if the user is reconnecting tell him that the server has restarted
			if message.Reconnecting {
				return impl.StructMessageResponse(responses.UserStateRegeneration{
					State: responses.ServerRestarted,
					User:  *user,
				})
			} else {
				return impl.StructMessageResponse(resp)
			}
		}
	}

	if message.Reconnecting {
		return impl.StructMessageResponse(responses.UserStateRegeneration{
			State: responses.ServerRestartedNameTaken,
			User:  model.User{},
		})
	} else {
		return impl.ErrorResponse(fmt.Sprintf("User of name %s is already logged on the server under ID of %d", user.Name, user.ID), impl.PlayerNameAlreadyTaken)
	}
}

// handler of OnClientDisconnected events
func (k *KrisKrosServer) OnClientDisconnected(socket int, leaving bool) {
	userID, exists := k.Router.SocketToUserID[socket]
	log.Warnf("User of Socket ID=%d has disconnected\n", socket)
	if exists {
		log.Warnf("Deleting a socket=%d and userID=%d from UserIDToSocket map", socket, userID)
		delete(k.Router.SocketToUserID, socket)
		delete(k.Router.UserIDToSocket, userID)
		user, exists := k.usersById[userID]
		if exists {
			state, exists := k.Router.UserStates[user.ID]
			if exists {
				if state.Id() >= LOBBY_JOINED_ID && state.Id() <= LOBBY_CREATED_READY_ID {
					k.removePlayerFromLobby(user.ID)
					log.Infof("Deleting a player of name=%s from lobby", user.Name)
				} else if state.Id() >= GAME_STARTED_STATE_ID && state.Id() <= WORDS_VALIDITY_DECIDED_STATE_ID {
					k.gameServer.PlayerLeft(userID, state.Id(), leaving)
				} else {
					// TODO
				}
			}
		}
	}
}

// handler of OnUserDisconnecting events, a willful disconnect
func (k *KrisKrosServer) OnUserDisconnecting(userID int) {
	user, exists := k.usersById[userID]
	if exists {
		socket, exists := k.Router.UserIDToSocket[userID]
		if exists {
			k.OnClientDisconnected(socket, true)
			delete(k.usersById, userID)
			delete(k.usersByName, user.Name)
			delete(k.Router.UserStates, user.ID)
		}
	}
}

func (k *KrisKrosServer) onKeepAlive(userID int) def.Response {
	return impl.StructMessageResponse(responses.KeepAliveResponse{})
}
