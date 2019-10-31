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
		1:   CreatedMessageHandler{},
		"á": 2,
		"b": 2,
		"c": 3,
		"č": 4,

		"d": 2,
		"ď": 8,
		"e": 1,
		"é": 5,
		"ě": 5,

		"f":  8,
		"g":  8,
		"h":  3,
		"ch": 4,
		"i":  1,

		"í": 2,
		"j": 2,
		"k": 1,
		"l": 1,
		"m": 2,

		"n": 1,
		"ň": 6,
		"o": 1,
		"ó": 10,
		"p": 1,

		"r": 1,
		"ř": 4,
		"s": 1,
		"š": 3,
		"t": 1,

		"ť": 6,
		"u": 2,
		"ů": 5,
		"ú": 6,
		"v": 1,

		"x": 10,
		"y": 1,
		"ý": 4,
		"z": 3,
		"ž": 4,
	}
}

func (s SimpleJsonReader) Read(message SimpleMessage) {

}

func (s SimpleJsonReader) SetOutput(reader ApplicationMessageReader) {
	panic("implement me")
}
