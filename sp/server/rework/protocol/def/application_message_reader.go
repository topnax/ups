package def

type ApplicationMessageReader interface {
	Read(message MessageHandler, clientUID int) Response
}
