package responses

import (
	"ups/sp/server/game"
	"ups/sp/server/model"
)

const (
	getLobbiesResponse      = 101
	playerJoinedResponse    = 102
	lobbyJoinedResponse     = 103
	playerLeftLobbyResponse = 104
	lobbyDestroyedResponse  = 105
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
	Lobbies []model.Lobby `json:"lobbies"`
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

//////////////////////////////////////

type LobbyJoinedResponse struct {
	Lobby model.Lobby `json:"lobby"`
}

func (g LobbyJoinedResponse) Type() int {
	return lobbyJoinedResponse
}

//////////////////////////////////////

type PlayerLeftLobby struct {
	Player game.Player `json:"player"`
}

func (g PlayerLeftLobby) Type() int {
	return playerLeftLobbyResponse
}

//////////////////////////////////////

type LobbyDestroyed struct {
}

func (g LobbyDestroyed) Type() int {
	return lobbyDestroyedResponse
}
