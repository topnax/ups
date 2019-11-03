package encoding

type JsonReader interface {
	Read(message SimpleMessage)
	SetOutput(reader ApplicationMessageReader)
}

type SimpleJsonReader struct {
}

type MessageHandler interface {
	handle(message SimpleMessage)
}

type CreatedMessageHandler struct{}

func (c CreatedMessageHandler) handle(message SimpleMessage) {
	panic("implement me")
}

func GetMessageHandlers() map[int]MessageHandler {
	return map[int]MessageHandler{
		1: CreatedMessageHandler{},
	}
}

func (s SimpleJsonReader) Read(message SimpleMessage) {

}

func (s SimpleJsonReader) SetOutput(reader ApplicationMessageReader) {
	panic("implement me")
}
