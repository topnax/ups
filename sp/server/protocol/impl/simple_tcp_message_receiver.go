package impl

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode/utf8"
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
	buffer      string
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

	var messages []string
	var prevChar rune
	lastGroupStart := 0

	runes := []rune(message)

	for pos, char := range message {
		if char == StartChar && prevChar != '\\' && lastGroupStart != pos {
			messages = append(messages, string(runes[lastGroupStart:pos]))
			lastGroupStart = pos
		}
		prevChar = char
	}

	utfLen := utf8.RuneCountInString(message)
	if lastGroupStart != utfLen {
		messages = append(messages, string(runes[lastGroupStart:utfLen]))
	}

	for _, mess := range messages {
		s.ReceiveMessage(UID, mess)
	}
}

func (s *SimpleTcpMessageReceiver) ReceiveMessage(UID int, message string) {

	// if, after the removal of the line break, the message is empty, return
	if len(message) < 1 {
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

	buffer := s.buffers[UID]
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
			buffer.buffer = parts[3]
			s.checkBufferReady(buffer)
		}
	} else {
		buffer.buffer += message
		s.checkBufferReady(buffer)
	}
}

func (s *SimpleTcpMessageReceiver) checkBufferReady(buffer *SimpleTcpMessageBuffer) {
	strlen := utf8.RuneCountInString(buffer.buffer)
	//if len(buffer.buffer) == buffer.Length {
	if strlen == buffer.Length {
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
			content:     buffer.buffer,
			id:          buffer.MessageId,
			msgType:     buffer.MessageType,
		})
	} else {
		response = ErrorResponseID("Cannot send message to JSON parser because it's null", NoMessageReader, buffer.MessageId)
	}

	s.Send(response, buffer.ClientUID, buffer.MessageId)
	log.Debugln("Responding to client %d '%s'", buffer.ClientUID, response.Content())
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
		s.responseSender.Send(fmt.Sprintf("%c%d%s%d%s%d%s%s", StartChar, utf8.RuneCountInString(response.Content()), Separator, response.Type(), Separator, msgID, Separator, response.Content()), clientUID)
	} else {
		log.Errorln("Cannot send response because output is null")
	}
}

func (buffer *SimpleTcpMessageBuffer) reset() {
	buffer.buffer = ""
	buffer.Length = 0
	buffer.MessageType = 0
	buffer.MessageId = 0
}
