package impl

import (
	"encoding/json"
	"fmt"
)

const (
	ErrorPrefix = 400

	MarshalError     = 1
	NoMessageReader  = 2
	NoMessageHandler = 3
	FailedToParse    = 4
	FailedToCast     = 5
	FailedToRoute    = 6

	PlainSuccess = 701
)

func (s SimpleResponse) Content() string {
	return s.content
}

type SimpleResponse struct {
	content      string
	responseType int
}

func (s SimpleResponse) Type() int {
	return s.responseType
}

////////////////////////////////////////////////////

func GetResponse(content string, responseType int) SimpleResponse {
	return SimpleResponse{
		content:      content,
		responseType: responseType,
	}
}

func SuccessResponse(content string) SimpleResponse {
	return GetResponse(content, PlainSuccess)
}

func ErrorResponse(content string, errorType int) SimpleResponse {
	return GetResponse(content, errorType+ErrorPrefix)
}

func MessageResponse(message interface{}, messageType int) SimpleResponse {
	bytes, err := json.Marshal(message)
	if err == nil {
		return GetResponse(string(bytes), messageType)
	} else {
		return ErrorResponse(fmt.Sprintf("Could not marshal a message of type %d, error %s", messageType, err), MarshalError)
	}
}
