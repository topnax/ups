package game_server

import (
	log "github.com/sirupsen/logrus"
	"ups/sp/server/encoding"
	"ups/sp/server/game"
)

type KrisKrosServer struct {
	lobbyUIQ int

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

func (k *KrisKrosServer) OnJoinLobby(message encoding.JoinLobbyMessage, clientUID int) {
	log.Infof("Receiver join message from %d", clientUID)
	_, exists := k.lobbies[message.LobbyID]
	if exists {
		player := game.Player{
			Name: message.ClientName,
			ID:   clientUID,
		}

		lobby := k.lobbies[message.LobbyID]
		lobby.Players = append(lobby.Players, player)

		log.Debugf("Lobby #%d", lobby.ID)
		for _, player := range lobby.Players {
			log.Debugf("[%d] %s", player.ID, player.Name)
		}
	}
}

func (k *KrisKrosServer) OnCreateLobby(message encoding.CreateLobbyMessage, clientUID int) {
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

		log.Debugf("%s created new lobby of ID %d", player.Name, lobby.ID)
	} else {
		log.Errorf("Player #%d already created a lobby", clientUID)
	}
}

//
//func (k *KrisKrosServer) OnCreateLobby(message CreateLobbyMessage, clientUID int) {
//	//log.Infof("Receiver create message from %d, count %d", clientUID, k.count)
//	//k.count++
//}
//
//func (k KrisKrosServer) OnJoinLobby(message JoinLobbyMessage, clientUID int) {
//	log.Infof("Received join message to %d from %d", message.LobbyID, clientUID)
//}
