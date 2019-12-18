package protocol

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
)

type SimpleOutput struct {
	lastMessage def.Message
	messages    []def.Message
}

func (s *SimpleOutput) Read(message def.Message) def.Response {
	s.lastMessage = message
	s.messages = append(s.messages, message)
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

func TestReceive2(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	bytec := len([]byte("{\"player_name\" : \"alzáček\"}"))

	sent := fmt.Sprintf("$%d#2#3#{\"player_name\" : \"alzáček\"}", bytec)
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))
	if jsonReader.lastMessage.Content() != "{\"player_name\" : \"alzáček\"}" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "$%d#2#3#{\"player_name\" : \"alzáček\"}")
	} else if jsonReader.lastMessage.Type() != 2 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 2)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 3 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 3)
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

	if jsonReader.lastMessage.Content() != "Hello$guy" {
		t.Errorf("Got message %s, want %s", jsonReader.lastMessage.Content(), "Hello$guy")
	} else if jsonReader.lastMessage.Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.lastMessage.Type(), 5)
	} else if jsonReader.lastMessage.ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ClientID(), 1)
	} else if jsonReader.lastMessage.ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.lastMessage.ID(), 10)
	}
}

func TestReceiveMultiple(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	bytec := len([]byte("Pěšák"))

	sent := fmt.Sprintf("$10#5#10#Hello\\$guy$%d#7#11#Pěšák", bytec)
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.messages[0].Content() != "Hello$guy" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[0].Content(), "Hello$guy")
	} else if jsonReader.messages[0].Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[0].Type(), 5)
	} else if jsonReader.messages[0].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ClientID(), 1)
	} else if jsonReader.messages[0].ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ID(), 10)
	}

	if jsonReader.messages[1].Content() != "Pěšák" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[1].Content(), "Pěšák")
	} else if jsonReader.messages[1].Type() != 7 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[1].Type(), 7)
	} else if jsonReader.messages[1].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ClientID(), 1)
	} else if jsonReader.messages[1].ID() != 11 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ID(), 10)
	}

}

func TestReceiveMultiple2(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	bytec := len([]byte("Pěšáčečíkk"))

	sent := fmt.Sprintf("$10#5#10#Hello\\$guy$%d#7#11#Pěšáčečíkk", bytec)
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.messages[0].Content() != "Hello$guy" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[0].Content(), "Hello$guy")
	} else if jsonReader.messages[0].Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[0].Type(), 5)
	} else if jsonReader.messages[0].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ClientID(), 1)
	} else if jsonReader.messages[0].ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ID(), 10)
	}

	if jsonReader.messages[1].Content() != "Pěšáčečíkk" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[1].Content(), "Pěšák")
	} else if jsonReader.messages[1].Type() != 7 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[1].Type(), 7)
	} else if jsonReader.messages[1].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ClientID(), 1)
	} else if jsonReader.messages[1].ID() != 11 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ID(), 10)
	}

}

func TestReceiveMultipleSplit(t *testing.T) {
	jsonReader := SimpleOutput{}
	smr := impl.SimpleTcpMessageReceiver{}
	smr.SetMessageReader(&jsonReader)

	bytec := len([]byte("Pěšák"))

	sent := fmt.Sprintf("$10#5#10#Hello\\$guy$%d#7#11#Pěš", bytec)
	bytes := []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	sent = "ák"
	bytes = []byte(sent)
	smr.Receive(1, bytes, len(bytes))

	if jsonReader.messages[0].Content() != "Hello$guy" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[0].Content(), "Hello$guy")
	} else if jsonReader.messages[0].Type() != 5 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[0].Type(), 5)
	} else if jsonReader.messages[0].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ClientID(), 1)
	} else if jsonReader.messages[0].ID() != 10 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[0].ID(), 10)
	}

	if jsonReader.messages[1].Content() != "Pěšák" {
		t.Errorf("Got message %s, want %s", jsonReader.messages[1].Content(), "Pěšák")
	} else if jsonReader.messages[1].Type() != 7 {
		t.Errorf("Got message type %d, want %d", jsonReader.messages[1].Type(), 7)
	} else if jsonReader.messages[1].ClientID() != 1 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ClientID(), 1)
	} else if jsonReader.messages[1].ID() != 11 {
		t.Errorf("Got client UID %d, want %d", jsonReader.messages[1].ID(), 10)
	}

}
