package encoding

type ApplicationMessageReader interface {
	OnCreatedMessageReceived(message Crea)
}

type ApplicationMessage struct {
	message     interface{}
	messageType int
}
