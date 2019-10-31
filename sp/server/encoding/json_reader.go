package encoding

type JsonReader interface {
	Read(message Message)
	SetOutput(reader ApplicationMessageReader)
	//SetOutput(M)
}
