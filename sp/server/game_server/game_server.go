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

	count             int
	lobbies           map[int]*Lobby
	lobbiesByOwnerID  map[int]*Lobby
	lobbiesByPlayerID map[int]*Lobby
}

func (k *KrisKrosServer) ClientDisconnected(clientUID int) {
	k.removeClientFromLobby(clientUID)
	delete(k.lobbiesByPlayerID, clientUID)
}

func (k *KrisKrosServer) removeClientFromLobby(clientUID int) {
	lobby, ok := k.lobbiesByPlayerID[clientUID]

	if ok {
		playerName := ""
		for _, player := range lobby.Players {
			if player.ID == clientUID {
				playerName = player.Name
			}
		}

		playerIndex := -1
		for i, player := range lobby.Players {
			if player.ID != clientUID {
				k.SendMessage(encoding.PlayerLeftLobbyMessage{ClientName: playerName}, player.ID)
			} else {
				playerIndex = i
			}
		}
		lobby.Players = append(lobby.Players[:playerIndex], lobby.Players[playerIndex+1:]...)
	}
}

func NewKrisKrosServer() KrisKrosServer {
	return KrisKrosServer{
		lobbies:           make(map[int]*Lobby),
		lobbiesByOwnerID:  make(map[int]*Lobby),
		lobbiesByPlayerID: make(map[int]*Lobby),
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

		k.lobbiesByPlayerID[newPlayer.ID] = lobby

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
			if player.ID != newPlayer.ID {
				k.SendMessage(encoding.PlayerJoinedLobbyMessage{ClientName: newPlayer.Name}, player.ID)
			}
		}

		return encoding.SuccessResponse(fmt.Sprintf("Sucessfully joined a lobby of ID %d. Owner %s", lobby.ID, lobby.Owner.Name))
	}
	return encoding.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID))
}

func (k *KrisKrosServer) OnCreateLobby(message encoding.CreateLobbyMessage, clientUID int) encoding.ResponseMessage {
	log.Infof("Receiver create message from %d", clientUID)
	_, exists := k.lobbiesByOwnerID[clientUID]
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

		k.lobbiesByOwnerID[clientUID] = k.lobbies[k.lobbyUIQ]
		lobby.Players = append(lobby.Players, player)

		return encoding.SuccessResponse(fmt.Sprintf("%s created new lobby of ID %d", player.Name, lobby.ID))
	} else {
		return encoding.ErrorResponse(fmt.Sprintf("Player #%d already created a lobby", clientUID))
	}
}
