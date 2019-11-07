package encoding

import (
	"encoding/json"
	"fmt"
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
}

func GetResponse(Type int, Content string) ResponseMessage {
	return ResponseMessage{
		Type:    Type,
		Content: Content,
	}
}

func SuccessResponse(Content string) ResponseMessage {
	return GetResponse(1, Content)
}

func ErrorResponse(Content string) ResponseMessage {
	return GetResponse(0, Content)
}

func MessageResponse(Struct interface{}, Type int) ResponseMessage {
	bytes, err := json.Marshal(Struct)
	if err == nil {
		return GetResponse(Type, string(bytes))
	} else {
		return ErrorResponse(fmt.Sprintf("Could not marshal a message of type %d, error %s", Type, err))
	}
}
