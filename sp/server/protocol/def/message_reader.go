package def

// an interface that defines a message reader, that reads a parsed TCP message with fields
type MessageReader interface {
	Read(message Message) Response
}

// defines a message handler, that takes a  parsed TCP message and returns a response. Each handler has to have a type
type MessageHandler interface {
	Handle(message Message, amr ApplicationMessageReader) Response
	GetType() int
}
