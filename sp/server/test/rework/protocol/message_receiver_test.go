package protocol

import (
	"testing"
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/impl"
)

type SimpleOutput struct {
	lastMessage def.Message
}

func (s SimpleOutput) Send(content string, clientUID int) {

}

func (s *SimpleOutput) Read(message def.Message) def.Response {
	println("received read")
	s.lastMessage = message
	return impl.SuccessResponse("test")
}
func TestReceive(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)
	sent := "$10#5#10#Hello guys"
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))
	if jsonReader.lastMessage.Content() != "Hello guys" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "Hello guys")
	} else if jsonReader.lastMessage.Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 5)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 10)
	}
}

func TestReceivePartial(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	sent := "$10#5#10#Hello"
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	sent = " guys"
	bytes = []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.lastMessage.Content() != "Hello guys" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "Hello guys")
	} else if jsonReader.lastMessage.Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 5)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 10)
	}
}

func TestReceivePartial2(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	sent := "$10#5#10#Youre"
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	sent = " guys hello"
	bytes = []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	sent = "$10#5#10#Hello guys"
	bytes = []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.lastMessage.Content() != "Hello guys" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "Hello guys")
	} else if jsonReader.lastMessage.Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 5)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 10)
	}
}

func TestReceivePartial3(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	sent := "$10#5#10#Hello\\"
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	sent = "$guy"
	bytes = []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.lastMessage.Content() != "Hello\\$guy" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "Hello\\$guy")
	} else if jsonReader.lastMessage.Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 5)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 10)
	}
}
