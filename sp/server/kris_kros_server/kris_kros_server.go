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
			Name: message.PlayerName,
			ID:   clientUID,
		}

		lobby := k.lobbies[message.LobbyID]
		lobby.Players = append(lobby.Players, newPlayer)

		k.lobbiesByPlayerID[newPlayer.ID] = lobby

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
			if player.ID != newPlayer.ID {
				resp := responses.PlayerJoinedResponse{PlayerName: newPlayer.Name, PlayerID: newPlayer.ID}
				k.sender.Send(impl.MessageResponse(resp, resp.Type()), player.ID, 0)
			}
		}

		return impl.SuccessResponse(fmt.Sprintf("Sucessfully joined a lobby of ID %d. Owner %s", lobby.ID, lobby.Owner.Name))
	}
	return impl.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID), impl.LobbyDoesNotExist)
}

func (k *KrisKrosServer) OnCreateLobby(msg messages.CreateLobbyMessage, clientUID int) def.Response {
	log.Infof("Receiver create message from %d", clientUID)
	_, exists := k.lobbiesByOwnerID[clientUID]
	if !exists {
		player := game.Player{
			Name: msg.PlayerName,
			ID:   clientUID,
		}
		k.lobbies[k.lobbyUIQ] = &model.Lobby{
			Owner: player,
			ID:    k.lobbyUIQ,
		}
		lobby := k.lobbies[k.lobbyUIQ]
		k.lobbyUIQ++

		k.lobbiesByOwnerID[clientUID] = k.lobbies[k.lobbyUIQ]
		lobby.Players = append(lobby.Players, player)

		return impl.SuccessResponse(fmt.Sprintf("lobby created OF ID %d", lobby.ID))
	} else {
		return impl.ErrorResponse(fmt.Sprintf("Player #%d already created a lobby", clientUID), impl.PlayerAlreadyCreatedLobby)
	}
}

func (k *KrisKrosServer) OnGetLobbies(message messages.GetLobbiesMessage, clientID int) def.Response {
	var lobbies []model.Lobby

	for _, lobby := range k.lobbies {
		lobbies = append(lobbies, *lobby)
	}

	resp := responses.GetLobbiesResponse{Lobbies: lobbies}

	return impl.MessageResponse(resp, resp.Type())
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
