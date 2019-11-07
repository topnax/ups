package encoding

type ApplicationMessageReader interface {
	SetMessageSender(sender MessageSender)
	OnJoinLobby(message JoinLobbyMessage, clientUID int) ResponseMessage
	OnCreateLobby(message CreateLobbyMessage, clientUID int) ResponseMessage
}
