package encoding

import (
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

// #{len}#{type}#{json}
// #15#1#{name:"Standa"}

const (
	SEPARATOR = "#"
)

type MessageReader interface {
	Receive(UID int, bytes []byte, length int)
	SetOutput(channel chan Message)
}

type SimpleMessageReader struct {
	messageCount int
	buffers      map[int]*SimpleMessageBuffer
	output       io.Reader
}

type SimpleMessageBuffer struct {
	Length      int
	MessageType int
	buffer      string
}

func (s *SimpleMessageReader) Receive(UID int, bytes []byte, length int) {

	message := string(bytes[:length])

	if last := len(message) - 1; last >= 0 && message[last] == '\n' {
		message = message[:last]
	}

	log.Infoln(len(message))

	for i, char := range message {
		log.Infof("[%d] '%s' %d", i, string(char), int(char))
	}

	//}
	log.Infof("received from #%d - '%s'", UID, message)

	if s.buffers == nil {
		log.Infoln("Buffer map not created yet... Creating new")
		s.buffers = make(map[int]*SimpleMessageBuffer)
	}

	_, exists := s.buffers[UID]

	if !exists {

		s.buffers[UID] = &SimpleMessageBuffer{
			Length: 0,
			buffer: "",
		}

	}

	buffer := s.buffers[UID]
	log.Infoln("Buffer len is", buffer.Length)
	if len(s.buffers[UID].buffer) <= 0 {
		parts := strings.Split(message, SEPARATOR)
		if len(parts) != 4 {
			log.Infoln("Invalid first message part")
			return
		}
		length, err := strconv.Atoi(parts[1])
		messageType, err2 := strconv.Atoi(parts[2])

		log.Infof("Message '%s'\n [1] '%s'\n [2] '%s'\n [3] '%s'\n", message, parts[1], parts[2], parts[3])

		if err == nil && err2 == nil {
			buffer.Length = length
			buffer.MessageType = messageType
			buffer.buffer += parts[3]
			if len(buffer.buffer) == length {
				log.Infof("Parsed from %d at first '%s'\n length %d\n type %d", UID, buffer.buffer, buffer.Length, buffer.MessageType)
				buffer.buffer = ""
				buffer.Length = 0
				buffer.MessageType = 0
			}
		}
	} else {
		log.Infof("Adding to buffer...")
		buffer.buffer += message
		log.Infof("Current buffer len is %d\n\n'%s'", len(buffer.buffer), buffer.buffer)
		if len(buffer.buffer) == buffer.Length {
			log.Infof("Parsed from %d later '%s'\n length %d\n type %d", UID, buffer.buffer, buffer.Length, buffer.MessageType)
			buffer.buffer = ""
			buffer.Length = 0
			buffer.MessageType = 0
		}
	}

	////messageContent := "Hello dudes"
	//if s.output != nil {
	//	//s.output.
	//	//	s.outChannel <- Message{
	//	//	Length:  len(messageContent),
	//	//	Type:    2,
	//	//	Content: []byte(messageContent),
	//	//}
	//} else {
	//	logrus.Errorln("cannot send message because out channel is null")
	//}
}

func (s *SimpleMessageReader) SetOutput(channel chan Message) {
	//s.outChannel = channel
}

type Message struct {
	ClientUID int
	Length    int
	Type      int
	Content   []byte
}
