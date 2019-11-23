package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/game"
	"ups/sp/server/model"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
	"ups/sp/server/protocol/responses"
)

type KrisKrosServer struct {
	count  int
	Router KrisKrosRouter
	sender def.ResponseSender

	lobbyUIQ int

	lobbies           map[int]*model.Lobby
	lobbiesByOwnerID  map[int]*model.Lobby
	lobbiesByPlayerID map[int]*model.Lobby
}

func failedToCast(message def.MessageHandler) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed to cast message of type %d", message.GetType()), impl.FailedToCast)
}

func NewKrisKrosServer(sender def.ResponseSender) KrisKrosServer {
	kks := KrisKrosServer{
		lobbies:           make(map[int]*model.Lobby),
		lobbiesByOwnerID:  make(map[int]*model.Lobby),
		lobbiesByPlayerID: make(map[int]*model.Lobby),
	}

	kks.Router = newKrisKrosRouter(&kks)
	kks.sender = sender
	return kks
}

func (k KrisKrosServer) Read(message def.MessageHandler, clientUID int) def.Response {
	res := k.Router.route(message, clientUID)
	log.Infoln("Server returning: ", res.Content())
	return res
}

func (k *KrisKrosServer) OnPlayerJoined(message messages.PlayerJoinedMessage, clientUID int) def.Response {
	log.Infoln("On player joined:)", message.PlayerName, k.count)
	k.count++
	return impl.SuccessResponse("Successfully :) " + message.PlayerName)
}

func (k *KrisKrosServer) OnJoinLobby(message messages.JoinLobbyMessage, clientUID int) def.Response {
	_, exists := k.lobbies[message.LobbyID]
	if exists {
		newPlayer := game.Player{
			Name:  message.PlayerName,
			ID:    clientUID,
			Ready: false,
		}

		lobby := k.lobbies[message.LobbyID]
		lobby.Players = append(lobby.Players, newPlayer)

		k.lobbiesByPlayerID[newPlayer.ID] = lobby

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
			if player.ID != newPlayer.ID {
				//resp := responses.PlayerJoinedResponse{PlayerName: newPlayer.Name, PlayerID: newPlayer.ID}
				resp := responses.LobbyJoinedResponse{Lobby: *lobby}
				k.sender.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
			}
		}
		resp := responses.LobbyJoinedResponse{Lobby: *lobby}
		return impl.MessageResponse(resp, resp.Type())
	}
	return impl.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID), impl.LobbyDoesNotExist)
}

func (k *KrisKrosServer) OnCreateLobby(msg messages.CreateLobbyMessage, clientUID int) def.Response {
	log.Infof("Receiver create message from %d", clientUID)
	_, exists := k.lobbiesByOwnerID[clientUID]
	if !exists {
		player := game.Player{
			Name:  msg.PlayerName,
			ID:    clientUID,
			Ready: false,
		}
		k.lobbies[k.lobbyUIQ] = &model.Lobby{
			Owner: player,
			ID:    k.lobbyUIQ,
		}
		lobby := k.lobbies[k.lobbyUIQ]
		k.lobbiesByPlayerID[player.ID] = lobby

		k.lobbyUIQ++

		k.lobbiesByOwnerID[clientUID] = k.lobbies[k.lobbyUIQ]
		lobby.Players = append(lobby.Players, player)

		return impl.SuccessResponse(fmt.Sprintf("lobby created OF ID %d", lobby.ID))
	} else {
		return impl.ErrorResponse(fmt.Sprintf("Player #%d already created a lobby", clientUID), impl.PlayerAlreadyCreatedLobby)
	}
}

func (k *KrisKrosServer) OnGetLobbies(message messages.GetLobbiesMessage, clientID int) def.Response {
	lobbies := []model.Lobby{}

	for _, lobby := range k.lobbies {
		lobbies = append(lobbies, *lobby)
	}

	resp := responses.GetLobbiesResponse{Lobbies: lobbies}

	return impl.MessageResponse(resp, resp.Type())
}

func (k *KrisKrosServer) ClientDisconnected(clientUID int) {
	k.removeClientFromLobby(clientUID)
}

func (k *KrisKrosServer) removeClientFromLobby(clientUID int) bool {
	lobby, ok := k.lobbiesByPlayerID[clientUID]
	log.Infof("Inside remove client from lobby, clientuid %d, exists : %v", clientUID, ok)
	if ok {
		var dcdPlayer game.Player
		dcdPlayerIndex := -1
		for i, player := range lobby.Players {
			if player.ID == clientUID {
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
			delete(k.lobbiesByPlayerID, clientUID)

			for _, player := range lobby.Players {
				if player.ID != clientUID {
					resp := responses.LobbyJoinedResponse{Lobby: *lobby}
					k.sender.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
				}
			}
		}
		return false
	}
	return false
}

func (k *KrisKrosServer) destroyLobby(lobby *model.Lobby) {
	log.Infof("Destroying a lobby of id %d and owner %s", lobby.ID, lobby.Owner.Name)
	for _, player := range lobby.Players {
		if player.ID != lobby.Owner.ID {
			resp := responses.LobbyDestroyed{}
			delete(k.lobbiesByPlayerID, player.ID)
			k.sender.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
		}
	}

	delete(k.lobbiesByOwnerID, lobby.Owner.ID)
	delete(k.lobbies, lobby.ID)
}

func (k *KrisKrosServer) OnLeaveLobby(clientUID int) def.Response {
	if k.removeClientFromLobby(clientUID) {
		return impl.SuccessResponse("Successfully left lobby")
	} else {
		return impl.ErrorResponse("Could not leave the lobby", impl.CouldNotLeaveLobby)
	}
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
					resp := responses.LobbyJoinedResponse{Lobby: *lobby}
					k.sender.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
				}
			}
			resp := responses.LobbyJoinedResponse{Lobby: *lobby}
			log.Infoln("Owner is", lobby.Owner.Ready)
			return impl.MessageResponse(resp, resp.Type())
		}
	}
	return impl.ErrorResponse("Could not find such user in a lobby", impl.CouldNotFindSuchUserInLobby)
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
