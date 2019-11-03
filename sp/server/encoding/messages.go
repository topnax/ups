package encoding

import "github.com/sirupsen/logrus"

func (s *SimpleJsonReader) Init() {
	s.handlers = GetMessageHandlers()
}

type CreatedMessageHandler struct {
	Surname string        `json:"surname"`
	Smr     SampleMessage `json:"smr"`
}

type SampleMessage struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c *CreatedMessageHandler) Handle(message SimpleMessage) {
	if message.Parse(&c) {
		logrus.Infof("CreatedMessageReceived, surname %s,name %s, age %d", c.Surname, c.Smr.Name, c.Smr.Age)
	}
}

func (s *SampleMessage) Handle(message SimpleMessage) {
	if message.Parse(&s) {
		logrus.Infof("Simple message received, name %s, age %d", s.Name, s.Age)
	}
}
