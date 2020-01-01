package def

// defines a tcp message receiver that is able to receive tcp messages from the tcp layer
type TcpMessageReceiver interface {
	Receive(socket int, bytes []byte, length int)
}

// defines a message that must have an ID, a TYPE a client ID (socket) that the response will eventually be sent to
type Message interface {
	ID() int
	Type() int
	ClientID() int
	Content() string
}

// defines a reponse that has to have a content, a type an an ID of the message it is responding to
type Response interface {
	Content() string
	Type() int
	ID() int
}

// defines a message sender that sends a  message to the given socket
type MessageSender interface {
	Send(content string, socket int)
}

// defines a response sender that sends a response to the given socket
type ResponseSender interface {
	Send(response Response, socket int, msgID int)
}
