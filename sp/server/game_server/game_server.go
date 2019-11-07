package game_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/encoding"
	"ups/sp/server/game"
)

type KrisKrosServer struct {
	messageSender encoding.MessageSender
	lobbyUIQ      int

	count         int
	lobbies       map[int]*Lobby
	lobbiesByUser map[int]*Lobby
}

func NewKrisKrosServer() KrisKrosServer {
	return KrisKrosServer{
		lobbies:       make(map[int]*Lobby),
		lobbiesByUser: make(map[int]*Lobby),
	}
}

func (k *KrisKrosServer) SetMessageSender(messageSender encoding.MessageSender) {
	k.messageSender = messageSender
}

func (k *KrisKrosServer) SendMessage(message encoding.TypedMessage, clientUID int) {
	if k.messageSender != nil {
		k.messageSender.Send(encoding.MessageResponse(message, message.GetType()), clientUID)
	}
}

func (k *KrisKrosServer) OnJoinLobby(message encoding.JoinLobbyMessage, clientUID int) encoding.ResponseMessage {
	log.Infof("Receiver join message from %d", clientUID)
	_, exists := k.lobbies[message.LobbyID]
	if exists {
		newPlayer := game.Player{
			Name: message.ClientName,
			ID:   clientUID,
		}

		lobby := k.lobbies[message.LobbyID]
		lobby.Players = append(lobby.Players, newPlayer)

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
			if player.ID != newPlayer.ID {
				k.SendMessage(encoding.PlayerJoinedMessage{ClientName: newPlayer.Name}, player.ID)
			}
		}

		return encoding.SuccessResponse(fmt.Sprintf("Sucessfully joined a lobby of ID %d. Owner %s", lobby.ID, lobby.Owner.Name))
	}
	return encoding.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID))
}

func (k *KrisKrosServer) OnCreateLobby(message encoding.CreateLobbyMessage, clientUID int) encoding.ResponseMessage {
	log.Infof("Receiver create message from %d", clientUID)
	_, exists := k.lobbiesByUser[clientUID]
	if !exists {
		player := game.Player{
			Name: message.ClientName,
			ID:   clientUID,
		}
		k.lobbies[k.lobbyUIQ] = &Lobby{
			Owner: player,
			ID:    k.lobbyUIQ,
		}
		lobby := k.lobbies[k.lobbyUIQ]
		k.lobbyUIQ++

		k.lobbiesByUser[clientUID] = k.lobbies[k.lobbyUIQ]
		lobby.Players = append(lobby.Players, player)

		return encoding.SuccessResponse(fmt.Sprintf("%s created new lobby of ID %d", player.Name, lobby.ID))
	} else {
		return encoding.ErrorResponse(fmt.Sprintf("Player #%d already created a lobby", clientUID))
	}
}
