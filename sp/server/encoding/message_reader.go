package encoding

import (
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

const (
	START_CHAR = '$'
	SEPARATOR  = "#"
)

type MessageReader interface {
	Receive(UID int, bytes []byte, length int)
	SetOutput(channel chan SimpleMessage)
}

type SimpleMessageReader struct {
	messageCount int
	buffers      map[int]*SimpleMessageBuffer
	output       io.Reader
}

type SimpleMessageBuffer struct {
	ClientUID   int
	Length      int
	MessageType int
	buffer      string
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
		if len(parts) != 3 {
			log.Errorf("Invalid message header. Received message was `%s`", message)
			return
		}

		// parse message type and content length
		length, err := strconv.Atoi(parts[0])
		messageType, err2 := strconv.Atoi(parts[1])

		if err == nil && err2 == nil {
			// set buffer properties and append the message
			buffer.Length = length
			buffer.MessageType = messageType
			buffer.buffer = parts[2]
			s.checkBufferReady(buffer)
		}
	} else {
		buffer.buffer += message
		log.Debugf("Current buffer len is %d\n\nTotal message: '%s'", len(buffer.buffer), buffer.buffer)
		s.checkBufferReady(buffer)
	}
}

func (s *SimpleMessageReader) SetOutput(channel chan SimpleMessage) {
	//s.outChannel = channel
}

func (s *SimpleMessageReader) checkBufferReady(buffer *SimpleMessageBuffer) bool {
	return len(buffer.buffer) == buffer.Length
}

func (buffer *SimpleMessageBuffer) reset() {
	buffer.buffer = ""
	buffer.Length = 0
	buffer.MessageType = 0
}

type SimpleMessage struct {
	ClientUID int
	Length    int
	Type      int
	Content   string
}
