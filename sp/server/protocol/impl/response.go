package impl

import (
	"encoding/json"
	"fmt"
	"ups/sp/server/protocol/responses"
)

const (
	NoID        = 0
	ErrorPrefix = 400

	MarshalError     = 1
	NoMessageReader  = 2
	NoMessageHandler = 3
	FailedToParse    = 4
	FailedToCast     = 5
	FailedToRoute    = 6

	PlayerAlreadyCreatedLobby = 20
	LobbyDoesNotExist         = 21

	PlainSuccess = 701
)

func (s SimpleResponse) Content() string {
	return s.content
}

type SimpleResponse struct {
	content      string
	responseType int
	id           int
}

func (s SimpleResponse) Type() int {
	return s.responseType
}

func (s SimpleResponse) ID() int {
	return s.id
}

////////////////////////////////////////////////////

func GetResponse(content string, responseType int, id int) SimpleResponse {
	return SimpleResponse{
		content:      content,
		responseType: responseType,
		id:           id,
	}
}

func SuccessResponseID(content string, id int) SimpleResponse {
	return MessageResponseID(responses.PlainResponse{Content: content}, PlainSuccess, id)
}

func ErrorResponseID(content string, errorType int, id int) SimpleResponse {
	return MessageResponseID(responses.PlainResponse{Content: content}, errorType+ErrorPrefix, id)
}

func SuccessResponse(content string) SimpleResponse {
	return SuccessResponseID(content, NoID)
}

func ErrorResponse(content string, errorType int) SimpleResponse {
	return ErrorResponseID(content, errorType+ErrorPrefix, NoID)
}

func MessageResponseID(message interface{}, messageType int, id int) SimpleResponse {
	bytes, err := json.Marshal(message)
	if err == nil {
		return GetResponse(string(bytes), messageType, id)
	} else {
		return ErrorResponseID(fmt.Sprintf("Could not marshal a message of type %d, error %s", messageType, err), MarshalError, id)
	}
}

func MessageResponse(message interface{}, messageType int) SimpleResponse {
	bytes, err := json.Marshal(message)
	if err == nil {
		return GetResponse(string(bytes), messageType, 0)
	} else {
		return ErrorResponseID(fmt.Sprintf("Could not marshal a message of type %d, error %s", messageType, err), MarshalError, 0)
	}
}
