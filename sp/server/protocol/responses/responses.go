package responses

import "ups/sp/server/model"

const (
	getLobbiesResponse   = 101
	playerJoinedResponse = 102
)

type TypedResponse interface {
	Type() int
}

//////////////////////////////////////

type PlainResponse struct {
	Content string `json:"content"`
}

//////////////////////////////////////

type GetLobbiesResponse struct {
	Lobbies []model.Lobby
}

func (g GetLobbiesResponse) Type() int {
	return getLobbiesResponse
}

//////////////////////////////////////

type PlayerJoinedResponse struct {
	PlayerName string `json:"player_name"`
	PlayerID   int    `json:"player_id"`
}

func (g PlayerJoinedResponse) Type() int {
	return playerJoinedResponse
}
