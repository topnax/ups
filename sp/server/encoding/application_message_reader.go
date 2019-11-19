package encoding

type ApplicationMessageReader interface {
	SetMessageSender(sender MessageSender)
	//OnJoinLobby(message messages.JoinLobbyMessage, clientUID int) ResponseMessage
	//OnCreateLobby(message messages.CreateLobbyMessage, clientUID int) ResponseMessage
	//OnGetLobbies(message messages.GetLobbiesMessage, clientUID int) ResponseMessage
	OnJoinLobby(message Message, clientUID int) ResponseMessage
	OnCreateLobby(message Message, clientUID int) ResponseMessage
	OnGetLobbies(message Message, clientUID int) ResponseMessage
}
