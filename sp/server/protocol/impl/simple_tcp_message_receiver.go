package impl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
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
	ClientUID   int
	Length      int
	MessageType int
	buffer      []byte
	MessageId   int
}

type SimpleMessage struct {
	clientID    int
	messageType int
	content     string
	id          int
	msgType     int
}

func (s SimpleMessage) ID() int {
	return s.id
}

func (s SimpleMessage) Type() int {
	return s.msgType
}

func (s SimpleMessage) ClientID() int {
	return s.clientID
}

func (s SimpleMessage) Content() string {
	return s.content
}

// ADD A MESSAGE STARTING CHAR

func (s *SimpleTcpMessageReceiver) Receive(UID int, bytes []byte, length int) {

	// remove trailing empty bytes
	message := string(bytes[:length])

	// if, after the removal of the line break, the message is empty, return
	if len(message) < 1 {
		return
	}

	var messages [][]byte
	var prevChar rune
	lastGroupStart := 0

	for pos, char := range message {
		if char == StartChar && prevChar != '\\' && lastGroupStart != pos {
			messages = append(messages, bytes[lastGroupStart:pos])
			lastGroupStart = pos
		}
		prevChar = char
	}

	if lastGroupStart != len(bytes) {
		messages = append(messages, bytes[lastGroupStart:length])
	}

	for _, mess := range messages {
		log.Infof("Into receiving messages '%s'", mess)
		s.ReceiveMessage(UID, mess)
	}
}

func (s *SimpleTcpMessageReceiver) ReceiveMessage(UID int, bytes []byte) {

	// if, after the removal of the line break, the message is empty, return
	if len(bytes) < 1 {
		return
	}

	// check whether buffer map was created
	if s.buffers == nil {
		log.Debugln("Buffer map not created yet, creating new...")
		s.buffers = make(map[int]*SimpleTcpMessageBuffer)
	}

	_, exists := s.buffers[UID]

	if !exists {
		// no buffer was yet created for the given UID, so create a new one
		s.buffers[UID] = &SimpleTcpMessageBuffer{
			ClientUID: UID,
		}
	}
	message := string(bytes)
	buffer := s.buffers[UID]
	log.Debugf("Received message content is '%s'\n", message)

	if message[0] == StartChar && (len(s.buffers[UID].buffer) <= 0 || (len(buffer.buffer) > 0 && buffer.buffer[len(buffer.buffer)-1] != '\\')) {
		// if buffer length is equal or less than 0, a new message is received, empty the buffer
		parts := strings.Split(message[1:], Separator)
		if len(parts) != 4 {
			log.Errorf("Invalid message header. Received message was `%s`", message)
			return
		}

		// parse message type and content length
		length, err := strconv.Atoi(parts[0])
		messageType, err2 := strconv.Atoi(parts[1])
		messageId, err3 := strconv.Atoi(parts[2])

		if err == nil && err2 == nil && err3 == nil {
			// set buffer properties and append the message
			buffer.Length = length
			buffer.MessageType = messageType
			buffer.MessageId = messageId
			index := IndexOfNth(message, Separator, 3) + 1
			buffer.buffer = bytes[index:]
			s.checkBufferReady(buffer)
		}
	} else {
		//buffer.buffer += message
		buffer.buffer = append(buffer.buffer, bytes[0:]...)
		s.checkBufferReady(buffer)
	}
}

func (s *SimpleTcpMessageReceiver) checkBufferReady(buffer *SimpleTcpMessageBuffer) {
	//strlen := utf8.RuneCountInString(buffer.buffer)
	//if len(buffer.buffer) == buffer.Length {
	if len(buffer.buffer) == buffer.Length {
		s.clearBuffer(buffer)
	}
}

func (s *SimpleTcpMessageReceiver) clearBuffer(buffer *SimpleTcpMessageBuffer) {
	var response def.Response
	log.Infof("[#%d] %d - '%s'", buffer.ClientUID, buffer.MessageType, buffer.buffer)
	if s.messageReader != nil {
		response = s.messageReader.Read(SimpleMessage{
			clientID:    buffer.ClientUID,
			messageType: buffer.MessageType,
			content:     string(buffer.buffer),
			id:          buffer.MessageId,
			msgType:     buffer.MessageType,
		})
	} else {
		response = ErrorResponseID("Cannot send message to JSON parser because it's null", NoMessageReader, buffer.MessageId)
	}

	s.Send(response, buffer.ClientUID, buffer.MessageId)
	log.Debugf("Responding to client %d '%s'", buffer.ClientUID, response.Content())
	buffer.reset()
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

	log.Debugf("Sending message of type %d to %d: '%s'", response.Type(), clientUID, response.Content())
	if s.responseSender != nil {
		bytes := []byte(response.Content())
		s.responseSender.Send(fmt.Sprintf("%c%d%s%d%s%d%s%s", StartChar, len(bytes), Separator, response.Type(), Separator, msgID, Separator, response.Content()), clientUID)
	} else {
		log.Errorln("Cannot send response because output is null")
	}
}

func (buffer *SimpleTcpMessageBuffer) reset() {
	buffer.buffer = []byte{}
	buffer.Length = 0
	buffer.MessageType = 0
	buffer.MessageId = 0
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
