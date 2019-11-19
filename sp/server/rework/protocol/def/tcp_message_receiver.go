package def

type TcpMessageReceiver interface {
	Receive(UID int, bytes []byte, length int)
}

type Message interface {
	ID() int
	Type() int
	ClientID() int
	Content() string
}

type Response interface {
	Content() string
	Type() int
}

type MessageSender interface {
	Send(content string, clientUID int)
}

type ResponseSender interface {
	Send(response Response, clientUID int)
}
