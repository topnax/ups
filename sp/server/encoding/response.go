package encoding

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type ResponseOutput interface {
	Send(content string, clientUID int)
}

type Response struct {
	Content string `json:"content"`
}

type ResponseMessage struct {
	Type    int
	Content string
	ID      int
}

func GetResponse(Type int, Content string, ID int) ResponseMessage {
	logrus.Infoln("Creating a response")
	return ResponseMessage{
		Type:    Type,
		Content: Content,
		ID:      ID,
	}
}

func SuccessResponse(Content string) ResponseMessage {
	return GetResponse(100, Content, 0)
}

func ErrorResponse(Content string) ResponseMessage {
	return GetResponse(101, Content, 0)
}

func MessageResponse(Struct interface{}, Type int) ResponseMessage {
	bytes, err := json.Marshal(Struct)
	if err == nil {
		return GetResponse(Type, string(bytes), 0)
	} else {
		return ErrorResponse(fmt.Sprintf("Could not marshal a message of type %d, error %s", Type, err))
	}
}
