package encoding

type ApplicationMessageReader interface {
	Read(message ApplicationMessage)
}

type ApplicationMessage struct {
	message     interface{}
	messageType int
}
