package def

// an interface that defines a method that reads an application message and returns a response
type ApplicationMessageReader interface {
	Read(message MessageHandler, socket int) Response
}
