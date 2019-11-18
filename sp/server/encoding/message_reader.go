package encoding

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	START_CHAR = '$'
	SEPARATOR  = "#"
)

type MessageDecoder interface {
	Receive(UID int, bytes []byte, length int)
	SetOutput(jsonReader JsonReader)
}

type MessageSender interface {
	Send(response ResponseMessage, clientUID int)
}

type SimpleMessageReader struct {
	messageCount  int
	buffers       map[int]*SimpleMessageBuffer
	jsonReader    JsonReader
	networkOutput ResponseOutput
}

type SimpleMessage struct {
	ClientUID int
	Length    int
	Type      int
	Content   string
	ID        int
}

type SimpleMessageBuffer struct {
	ClientUID   int
	Length      int
	MessageType int
	buffer      string
	MessageId   int
}

// ADD A MESSAGE STARTING CHAR

func (s *SimpleMessageReader) Receive(UID int, bytes []byte, length int) {

	// remove trailing empty bytes
	message := string(bytes[:length])

	// if remove line break if found at the last position
	if last := len(message) - 1; last >= 0 && message[last] == '\n' {
		message = message[:last]
	}

	// if, after the removal of the line break, the message is empty, return
	if len(message) < 1 {
		return
	}

	log.Debugln("Received message from #%d => '%s'", UID, message)

	// check whether buffer map was created
	if s.buffers == nil {
		log.Debugln("Buffer map not created yet, creating new...")
		s.buffers = make(map[int]*SimpleMessageBuffer)
	}

	_, exists := s.buffers[UID]

	if !exists {
		// no buffer was yet created for the given UID, so create a new one
		s.buffers[UID] = &SimpleMessageBuffer{
			ClientUID: UID,
		}
	}

	buffer := s.buffers[UID]
	if message[0] == START_CHAR && (len(s.buffers[UID].buffer) <= 0 || (len(buffer.buffer) > 0 && buffer.buffer[len(buffer.buffer)-1] != '\\')) {
		// if buffer length is equal or less than 0, a new message is received, empty the buffer
		parts := strings.Split(message[1:], SEPARATOR)
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
		log.Debugf("Current buffer len is %d\n\nTotal message: '%s'", len(buffer.buffer), buffer.buffer)
		s.checkBufferReady(buffer)
	}
}

func (s *SimpleMessageReader) SetOutput(jsonReader JsonReader) {
	s.jsonReader = jsonReader
}

func (s *SimpleMessageReader) checkBufferReady(buffer *SimpleMessageBuffer) {
	if len(buffer.buffer) == buffer.Length {
		s.clearBuffer(buffer)
	}
}

func (s *SimpleMessageReader) clearBuffer(buffer *SimpleMessageBuffer) {
	var response ResponseMessage
	log.Infof("[#%d] %d - '%s'", buffer.ClientUID, buffer.MessageType, buffer.buffer)
	if s.jsonReader != nil {
		response = s.jsonReader.Read(SimpleMessage{
			ClientUID: buffer.ClientUID,
			Length:    buffer.Length,
			Type:      buffer.MessageType,
			Content:   buffer.buffer,
			ID:        buffer.MessageId,
		})
	} else {
		response = ErrorResponse("Cannot send message to JSON parser because it's null")
	}
	response.ID = buffer.MessageId
	buffer.reset()
	s.Send(response, buffer.ClientUID)
}

func (reader *SimpleMessageReader) SetResponseOutput(output ResponseOutput) {
	reader.networkOutput = output
}

func (s *SimpleMessageReader) Send(response ResponseMessage, clientUID int) {
	log.Debugf("Sending message of type %d to %d: '%s'", response.Type, clientUID, response.Content)
	if s.networkOutput != nil {
		log.Infof("%c%d%s%d%s%s", START_CHAR, len(response.Content), SEPARATOR, response.Type, SEPARATOR, response.Content)
		s.networkOutput.Send(fmt.Sprintf("%c%d%s%d%s%d%s%s", START_CHAR, len(response.Content), SEPARATOR, response.Type, SEPARATOR, response.ID, SEPARATOR, response.Content), clientUID)
	} else {
		log.Errorln("Cannot send response because output is null")
	}
}

func (buffer *SimpleMessageBuffer) reset() {
	buffer.buffer = ""
	buffer.Length = 0
	buffer.MessageType = 0
	buffer.MessageId = 0
}
