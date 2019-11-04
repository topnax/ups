package encoding

type ApplicationMessageReader interface {
	//OnCreateLobby(message CreateLobbyMessage)
	OnJoinLobby(message JoinLobbyMessage, clientUID int)
	OnCreateLobby(message CreateLobbyMessage, clientUID int)
}
