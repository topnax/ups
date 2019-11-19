package def

type MessageReader interface {
	Read(message Message) Response
}

type MessageHandler interface {
	Handle(message Message, amr ApplicationMessageReader) Response
	GetType() int
}
