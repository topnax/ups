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
	StartChar               = '$'
	Separator               = "#"
	InvalidByteCountCeiling = 5
)

type SimpleTcpMessageReceiver struct {
	messageCount   int
	buffers        map[int]*SimpleTcpMessageBuffer
	messageReader  def.MessageReader
	responseSender def.MessageSender
	socketCloser   def.SocketCloser
	TestMode       bool
}

type SimpleTcpMessageBuffer struct {
	State            int
	ClientUID        int
	Length           int
	MessageType      int
	headerBuffer     []byte
	contentBuffer    []byte
	MessageId        int
	InvalidByteCount int
	Closed           bool
}

// signals that an invalid byte has been receiver. If N invalid bytes are received the connection to the user is closed
func (buffer *SimpleTcpMessageBuffer) invalidByte(receiver *SimpleTcpMessageReceiver) {
	if !receiver.TestMode {
		if !buffer.Closed {
			buffer.InvalidByteCount++
			log.Warnf("Buffer of FD=%d received an invalid byte", buffer.ClientUID)
			if buffer.InvalidByteCount >= InvalidByteCountCeiling {
				buffer.Closed = true
				delete(receiver.buffers, buffer.ClientUID)
				log.Errorf("Connection to FD=%d is being closed because it's buffer received %d invalid bytes.", buffer.ClientUID, buffer.InvalidByteCount)
				if receiver.socketCloser != nil {
					receiver.socketCloser.CloseFd(buffer.ClientUID)
				}
			}
		}
	}
}

// signals that a valid byte has been received and that invalid byte count should be reset
func (buffer *SimpleTcpMessageBuffer) validByte() {
	buffer.InvalidByteCount = 0
}

// defines a simple message
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

// receives bytes from the given socket
func (s *SimpleTcpMessageReceiver) Receive(socket int, bytes []byte, length int) {

	// check whether headerBuffer map was created
	if s.buffers == nil {
		s.buffers = make(map[int]*SimpleTcpMessageBuffer)
	}

	_, exists := s.buffers[socket]

	if !exists {
		// no headerBuffer was yet created for the given socket, so create a new one
		s.buffers[socket] = &SimpleTcpMessageBuffer{
			ClientUID: socket,
			State:     1,
		}
	}
	message := string(bytes)
	buffer := s.buffers[socket]
	log.Debugf("TCP receiver received '%s'\n", message)

	// handling of received bytes is done using an automaton
	// for each socket an automaton is created
	// format of accepted messages is as following
	// {START_CHARACTER}{LENGTH}{SEPARATOR_CHARACTER}{MESSAGE_TYPE}{SEPARATOR_CHARACTER}{MESSAGE_ID}{SEPARATOR_CHARACTER}{CONTENT_OF_LENGTH_<LENGTH>}
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
				buffer.validByte()
			} else {
				buffer.invalidByte(s)
			}
		case 2:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 3
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 3:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 3
				buffer.validByte()
			} else if string(b) == Separator {
				length, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.Length = length
				buffer.State = 4
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 4:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 5
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 5:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 5
				buffer.validByte()
			} else if string(b) == Separator {
				messageType, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.MessageType = messageType
				buffer.State = 6
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 6:
			buffer.headerBuffer = []byte{}
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 7
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 7:
			if unicode.IsDigit(char) {
				buffer.headerBuffer = append(buffer.headerBuffer, b)
				buffer.State = 7
				buffer.validByte()
			} else if string(b) == Separator {
				messageId, _ := strconv.Atoi(string(buffer.headerBuffer))
				buffer.MessageId = messageId
				buffer.contentBuffer = []byte{}
				buffer.State = 8
				buffer.validByte()
			} else if char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.State = 1
				buffer.invalidByte(s)
			}
		case 8:
			if !IsNextByteEscaped(buffer.contentBuffer) && char == StartChar {
				buffer.State = 2
				buffer.invalidByte(s)
			} else {
				buffer.contentBuffer = append(buffer.contentBuffer, b)
				if len(buffer.contentBuffer) == buffer.Length {
					buffer.validByte()
					s.clearBuffer(buffer)
				} else if len(buffer.contentBuffer) > buffer.Length {
					buffer.invalidByte(s)
				} else {
					buffer.State = 8
					buffer.validByte()
				}
			}
		}
	}
}

// clears a buffer by sending the message to the message receiver and clearing it
func (s *SimpleTcpMessageReceiver) clearBuffer(buffer *SimpleTcpMessageBuffer) {

	if buffer.Closed {
		log.Warnf("Buffer of FD=%d not cleared, because it got closed")
		return
	}

	logLevel := log.GetLevel()
	if buffer.MessageType == 15 {
		// keep alive's logs are ignored
		log.SetLevel(log.WarnLevel)
	}
	var response def.Response
	log.Infof("[#%d] %d - '%s'", buffer.ClientUID, buffer.MessageType, buffer.contentBuffer)
	if s.messageReader != nil {
		// for test purposes do not escape \ character
		if !s.TestMode {
			response = s.messageReader.Read(SimpleMessage{
				clientID:    buffer.ClientUID,
				messageType: buffer.MessageType,
				content:     strings.Replace(string(buffer.contentBuffer), "\\", "", -1),
				id:          buffer.MessageId,
			})
		} else {
			response = s.messageReader.Read(SimpleMessage{
				clientID:    buffer.ClientUID,
				messageType: buffer.MessageType,
				content:     string(buffer.contentBuffer),
				id:          buffer.MessageId,
			})
		}
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

	if buffer.MessageType == 15 {
		log.SetLevel(logLevel)
	}
}

// returns true if the next byte would be escaped or not
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

func (receiver *SimpleTcpMessageReceiver) SetSocketCloser(remover def.SocketCloser) {
	receiver.socketCloser = remover
}

func (receiver *SimpleTcpMessageReceiver) SetOutput(output def.MessageSender) {
	receiver.responseSender = output
}

// sends a response in the same format as the one it accepts messages
func (s *SimpleTcpMessageReceiver) Send(response def.Response, clientUID int, msgID int) {

	if response.ID() != 0 {
		msgID = response.ID()
	}

	log.Debugf("About to send response of type %d to %d: '%s'", response.Type(), clientUID, response.Content())
	rawsponse := strings.Replace(response.Content(), Separator, "\\"+Separator, -1)
	//log.Debugf("First escapation '%s'", rawsponse)
	rawsponse = strings.Replace(rawsponse, string(StartChar), "\\"+string(StartChar), -1)
	//log.Debugf("Second escapation '%s'", rawsponse)

	if s.responseSender != nil {
		bytes := []byte(rawsponse)
		s.responseSender.Send(fmt.Sprintf("%c%d%s%d%s%d%s%s", StartChar, len(bytes), Separator, response.Type(), Separator, msgID, Separator, rawsponse), clientUID)
	} else {
		log.Errorln("Cannot send response because output is null")
	}
}

// returns the nth index of a string in the given string
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
