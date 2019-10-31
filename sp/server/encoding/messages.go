package encoding

import "encoding/json"

type SampleMessage struct {
	name string
	age  int
}

func (simpleMessage SimpleMessage) dos(messageType interface{}) {
	json.Unmarshal([]byte(simpleMessage.Content), messageType)

}

func (s SampleMessage) handle(message SimpleMessage) {
	message := &SimpleMessage{}
}
