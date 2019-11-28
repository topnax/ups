package responses

import (
	"ups/sp/server/game"
	"ups/sp/server/model"
)

const (
	ValidResponseCeiling = 400

	getLobbiesResponse        = 101
	playerJoinedResponse      = 102
	lobbyUpdatedResponse      = 103
	playerLeftLobbyResponse   = 104
	lobbyDestroyedResponse    = 105
	lobbyJoinedResponse       = 106
	userAuthenticatedResponse = 107
	lobbyStartedResponse      = 108
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

type LobbyUpdatedResponse struct {
	Lobby model.Lobby `json:"lobby"`
}

func (g LobbyUpdatedResponse) Type() int {
	return lobbyUpdatedResponse
}

//////////////////////////////////////

type PlayerLeftLobbyResponse struct {
	Player game.Player `json:"player"`
}

func (g PlayerLeftLobbyResponse) Type() int {
	return playerLeftLobbyResponse
}

//////////////////////////////////////

type LobbyDestroyedResponse struct {
}

func (g LobbyDestroyedResponse) Type() int {
	return lobbyDestroyedResponse
}

//////////////////////////////////////

type LobbyJoinedResponse struct {
	Player game.Player `json:"player"`
	Lobby  model.Lobby `json:"lobby"`
}

func (g LobbyJoinedResponse) Type() int {
	return lobbyJoinedResponse
}

//////////////////////////////////////

type UserAuthenticatedResponse struct {
	User model.User `json:"user"`
}

func (g UserAuthenticatedResponse) Type() int {
	return userAuthenticatedResponse
}

//////////////////////////////////////

type LobbyStartedResponse struct{}

func (g LobbyStartedResponse) Type() int {
	return lobbyStartedResponse
}
