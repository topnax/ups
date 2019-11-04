package encoding

import "github.com/sirupsen/logrus"

type ApplicationMessageReader interface {
	//OnCreateLobby(message CreateLobbyMessage)
	OnJoinLobby(message JoinLobbyMessage, clientUID int)
	OnCreateLobby(message CreateLobbyMessage, clientUID int)
}

type KrisKrosMessageReader struct {
	count int
}

func (k *KrisKrosMessageReader) OnCreateLobby(message CreateLobbyMessage, clientUID int) {
	logrus.Infof("Receiver create message from %d, count %d", clientUID, k.count)
	k.count++
}

func (k KrisKrosMessageReader) OnJoinLobby(message JoinLobbyMessage, clientUID int) {
	logrus.Infof("Received join message to %d from %d", message.LobbyID, clientUID)
}
