package kris_kros_server

//
//import (
//	"fmt"
//	log "github.com/sirupsen/logrus"
//	"ups/sp/server/encoding"
//	"ups/sp/server/game"
//	"ups/sp/server/messages"
//	"ups/sp/server/model"
//)
//
//type KrisKrosServer struct {
//	messageSender encoding.MessageSender
//	lobbyUIQ      int
//
//	count             int
//	lobbies           map[int]*model.Lobby
//	lobbiesByOwnerID  map[int]*model.Lobby
//	lobbiesByPlayerID map[int]*model.Lobby
//}
//
//func (k *KrisKrosServer) ClientDisconnected(clientUID int) {
//	k.removeClientFromLobby(clientUID)
//}
//
//func (k *KrisKrosServer) removeClientFromLobby(clientUID int) {
//	lobby, ok := k.lobbiesByPlayerID[clientUID]
//	log.Info("Inside remove client from lobby", clientUID, ok)
//	if ok {
//		playerName := ""
//		for _, player := range lobby.Players {
//			if player.ID == clientUID {
//				playerName = player.PlayerName
//			}
//		}
//
//		playerIndex := -1
//		for i, player := range lobby.Players {
//			if player.ID != clientUID {
//				k.SendMessage(messages.PlayerLeftLobbyMessage{ClientName: playerName}, player.ID)
//			} else {
//				playerIndex = i
//			}
//		}
//		lobby.Players = append(lobby.Players[:playerIndex], lobby.Players[playerIndex+1:]...)
//		delete(k.lobbiesByPlayerID, clientUID)
//	}
//}
//
//func NewKrisKrosServer() KrisKrosServer {
//	k := KrisKrosServer{
//		lobbies:           make(map[int]*model.Lobby),
//		lobbiesByOwnerID:  make(map[int]*model.Lobby),
//		lobbiesByPlayerID: make(map[int]*model.Lobby),
//	}
//
//	k.lobbies[1] = &model.Lobby{
//		ID:      7,
//		Players: nil,
//		Owner:   game.Player{},
//	}
//
//	return k
//}
//
//func (k *KrisKrosServer) SetMessageSender(messageSender encoding.MessageSender) {
//	k.messageSender = messageSender
//}
//
//func (k *KrisKrosServer) SendMessage(message encoding.TypedMessage, clientUID int) {
//	if k.messageSender != nil {
//		k.messageSender.Send(encoding.MessageResponse(message, message.GetType()), clientUID)
//	}
//}
//
//func (k *KrisKrosServer) OnJoinLobby(mes encoding.Message, clientUID int) encoding.ResponseMessage {
//	message := mes.(*messages.JoinLobbyMessage)
//	log.Infof("Receiver join message from %d", clientUID)
//	_, exists := k.lobbies[message.LobbyID]
//	if exists {
//		newPlayer := game.Player{
//			PlayerName: message.ClientName,
//			ID:   clientUID,
//		}
//
//		lobby := k.lobbies[message.LobbyID]
//		lobby.Players = append(lobby.Players, newPlayer)
//
//		k.lobbiesByPlayerID[newPlayer.ID] = lobby
//
//		log.Debugf("Lobby #%d", lobby.ID)
//		for _, player := range lobby.Players {
//			log.Debugf("[%d] %s", player.ID, player.PlayerName)
//			if player.ID != newPlayer.ID {
//				k.SendMessage(messages.PlayerJoinedLobbyMessage{ClientName: newPlayer.PlayerName}, player.ID)
//			}
//		}
//
//		return encoding.SuccessResponse(fmt.Sprintf("Sucessfully joined a lobby of ID %d. Owner %s", lobby.ID, lobby.Owner.PlayerName))
//	}
//	return encoding.ErrorResponse(fmt.Sprintf("Lobby of ID %d does not exist", message.LobbyID))
//}
//
//func (k *KrisKrosServer) OnCreateLobby(message encoding.Message, clientUID int) encoding.ResponseMessage {
//	log.Infof("Receiver create message from %d", clientUID)
//	_, exists := k.lobbiesByOwnerID[clientUID]
//	if !exists {
//		clmMes := message.(*messages.CreateLobbyMessage)
//		player := game.Player{
//			PlayerName: clmMes.ClientName,
//			ID:   clientUID,
//		}
//		k.lobbies[k.lobbyUIQ] = &model.Lobby{
//			Owner: player,
//			ID:    k.lobbyUIQ,
//		}
//		lobby := k.lobbies[k.lobbyUIQ]
//		k.lobbyUIQ++
//
//		k.lobbiesByOwnerID[clientUID] = k.lobbies[k.lobbyUIQ]
//		lobby.Players = append(lobby.Players, player)
//
//		return encoding.MessageResponse(messages.PlayerJoinedLobbyMessage{ClientName: "you hjoined"}, messages.PlayerJoinedLobbyMessage{}.GetType())
//	} else {
//		return encoding.ErrorResponse(fmt.Sprintf("Player #%d already created a lobby", clientUID))
//	}
//}
//
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
