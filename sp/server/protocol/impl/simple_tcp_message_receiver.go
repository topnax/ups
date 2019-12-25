package impl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode"
	"ups/sp/server/protocol/def"
)

const (
	StartChar = '$'
	Separator = "#"
)

type SimpleTcpMessageReceiver struct {
	messageCount   int
	buffers        map[int]*SimpleTcpMessageBuffer
	messageReader  def.MessageReader
	responseSender def.MessageSender
}

type SimpleTcpMessageBuffer struct {
	State         int
	ClientUID     int
	Length        int
	MessageType   int
	headerBuffer  []byte
	contentBuffer []byte
	MessageId     int
}

type SimpleMessage struct {
	clientID    int
	content     string
	id          int
	messageType int
}

func (s SimpleMessage) ID() int {
	return s.id
}

func (s SimpleMessage) Type() int {
	return s.messageType
}

func (s SimpleMessage) ClientID() int {
	return s.clientID
}

func (s SimpleMessage) Content() string {
	return s.content
}

// ADD A MESSAGE STARTING CHAR

func (s *SimpleTcpMessageReceiver) Receive(UID int, bytes []byte, length int) {

	// check whether headerBuffer map was created
	if s.buffers == nil {
		log.Debugln("Buffer map not created yet, creating new...")
		s.buffers = make(map[int]*SimpleTcpMessageBuffer)
	}

	_, exists := s.buffers[UID]

	if !exists {
		// no headerBuffer was yet created for the given UID, so create a new one
		s.buffers[UID] = &SimpleTcpMessageBuffer{
			ClientUID: UID,
			State:     1,
		}
	}
	message := string(bytes)
	buffer := s.buffers[UID]
	log.Debugf("Received message content is '%s'\n", message)

	for index, b := range bytes {

		char := rune(b)

		if index >= length {
			return
		}
		currentChar := string(b)

		_ = currentChar

		switch buffer.State {
		case 1:
			if char == StartChar {
				buffer.State = 2
			}
		case 2:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 3
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 3:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 3
			} else if string(b) == Separator {
				length, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.Length = length
				buffer.State = 4
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 4:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 5
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 5:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 5
			} else if string(b) == Separator {
				messageType, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.MessageType = messageType
				buffer.State = 6
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 6:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 7
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 7:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 7
			} else if string(b) == Separator {
				messageId, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.MessageId = messageId
				buffer.contentBuffer = []byte{}
				buffer.State = 8
			} else if char == StartChar {
				buffer.State = 2
			} else {
				buffer.State = 1
			}
		case 8:
			if !IsNextByteEscaped(buffer.contentBuffer) && char == StartChar {
				buffer.State = 2
			} else {
				buffer.contentBuffer = append(buffer.contentBuffer, b)
				if len(buffer.contentBuffer) == buffer.Length {
					s.clearBuffer(buffer)
				} else {
					buffer.State = 8
				}
			}
		}
	}
}

func (s *SimpleTcpMessageReceiver) clearBuffer(buffer *SimpleTcpMessageBuffer) {
	var response def.Response
	log.Infof("[#%d] %d - '%s'", buffer.ClientUID, buffer.MessageType, buffer.contentBuffer)
	if s.messageReader != nil {
		response = s.messageReader.Read(SimpleMessage{
			clientID:    buffer.ClientUID,
			messageType: buffer.MessageType,
			content:     strings.Replace(string(buffer.contentBuffer), "\\", "", -1),
			id:          buffer.MessageId,
		})
	} else {
		response = ErrorResponseID("Cannot send message to JSON parser because it's null", NoMessageReader, buffer.MessageId)
	}

	if response.Type() > -1 {
		s.Send(response, buffer.ClientUID, buffer.MessageId)
		log.Debugf("Responding to client %d '%s'", buffer.ClientUID, response.Content())
	} else {
		log.Debugf("Not responding to message of ID %d and Type %d", buffer.MessageId, buffer.MessageType)
	}

	buffer.State = 1
}

func IsNextByteEscaped(bytes []byte) bool {
	index := 0
	escCount := 0
	for len(bytes) > 0 && len(bytes) > index {
		if bytes[len(bytes)-index-1] == '\\' {
			escCount++
		} else {
			break
		}
		index++
	}
	return escCount%2 == 1
}

func (receiver *SimpleTcpMessageReceiver) SetMessageReader(reader def.MessageReader) {
	receiver.messageReader = reader
}

func (receiver *SimpleTcpMessageReceiver) SetOutput(output def.MessageSender) {
	receiver.responseSender = output
}

func (s *SimpleTcpMessageReceiver) Send(response def.Response, clientUID int, msgID int) {

	if response.ID() != 0 {
		msgID = response.ID()
	}

	log.Debugf("About to send response of type %d to %d: '%s'", response.Type(), clientUID, response.Content())
	rawsponse := strings.Replace(response.Content(), Separator, "\\"+Separator, -1)
	log.Debugf("First escapation '%s'", rawsponse)
	rawsponse = strings.Replace(rawsponse, string(StartChar), "\\"+string(StartChar), -1)
	log.Debugf("Second escapation '%s'", rawsponse)

	if s.responseSender != nil {
		bytes := []byte(rawsponse)
		s.responseSender.Send(fmt.Sprintf("%c%d%s%d%s%d%s%s", StartChar, len(bytes), Separator, response.Type(), Separator, msgID, Separator, rawsponse), clientUID)
	} else {
		log.Errorln("Cannot send response because output is null")
	}
}

func IndexOfNth(str string, tbf string, n int) int {
	if n < 1 {
		return -1
	}

	lastIndex := 0
	found := 0
	strlen := len(str)

	for {
		//fmt.Printf("substring from %s find %s at n=%d, lastIndex=%d, strlen=%d, substring=%s, indexof=%d\n", str, tbf, n, lastIndex, strlen, str[lastIndex:strlen],strings.Index(str[lastIndex:strlen], tbf))
		index := strings.Index(str[lastIndex:strlen], tbf)
		lastIndex += strings.Index(str[lastIndex:strlen], tbf)
		if index != -1 {
			found++
			if found == n {
				return lastIndex
			}
			lastIndex++
		} else {
			return -1
		}
	}
}
