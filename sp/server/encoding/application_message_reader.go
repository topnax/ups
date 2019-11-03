package encoding

import "github.com/sirupsen/logrus"

type ApplicationMessageReader interface {
	//OnCreateLobby(message CreateLobbyMessage)
	OnJoinLobby(message JoinLobbyMessage)
	OnCreateLobby(message CreateLobbyMessage)
}

type KrisKrosMessageReader struct {
	count int
}

func (k *KrisKrosMessageReader) OnCreateLobby(message CreateLobbyMessage, clientUID int) {
	logrus.Infof("Receiver create message from %d", clientUID)
	k.count++
}

func (k KrisKrosMessageReader) OnJoinLobby(message JoinLobbyMessage, clientUID int) {
	logrus.Infof("Received join message to %d from %d", message.LobbyID, clientUID)
}
