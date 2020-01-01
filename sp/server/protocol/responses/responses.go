package responses

import (
	"ups/sp/server/game"
	"ups/sp/server/model"
)

// DEFINITION OF RESPONSES USED

const (
	ValidResponseCeiling = 400

	getLobbiesResponse                = 101
	playerJoinedResponse              = 102
	lobbyUpdatedResponse              = 103
	playerLeftLobbyResponse           = 104
	lobbyDestroyedResponse            = 105
	lobbyJoinedResponse               = 106
	userAuthenticatedResponse         = 107
	lobbyStartedResponse              = 108
	gameStartedResponse               = 109
	tileUpdatedResponse               = 110
	tilesUpdatedResponse              = 111
	roundFinishedResponse             = 112
	playerAcceptedRoundResponse       = 113
	newRoundResponse                  = 114
	yourNewRoundResponse              = 115
	playerDeclinedWordsResponse       = 116
	gameEndedResponse                 = 117
	acceptResultedInNewRoundResponse  = 118
	playerConnectionChangedResponse   = 119
	gameStateRegenerationResponse     = 120
	keepAliveResponse                 = 121
	userStateRegenerationResponse     = 122
	finishResultedInNextRoundResponse = 123
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

//////////////////////////////////////

type GameStartedResponse struct {
	Players        []game.Player `json:"players"`
	Letters        []game.Letter `json:"letters"`
	ActivePlayerID int           `json:"active_player_id"`
}

func (g GameStartedResponse) Type() int {
	return gameStartedResponse
}

//////////////////////////////////////

type TileUpdatedResponse struct {
	Tile game.Tile `json:"tile"`
}

func (g TileUpdatedResponse) Type() int {
	return tileUpdatedResponse
}

//////////////////////////////////////

type TilesUpdatedResponse struct {
	Tiles                    []game.Tile `json:"tiles"`
	CurrentPlayerPoints      int         `json:"current_player_points"`
	CurrentPlayerTotalPoints int         `json:"current_player_total_points"`
}

func (g TilesUpdatedResponse) Type() int {
	return tilesUpdatedResponse
}

//////////////////////////////////////

type RoundFinishedResponse struct{}

func (g RoundFinishedResponse) Type() int {
	return roundFinishedResponse
}

//////////////////////////////////////

type PlayerAcceptedRoundResponse struct {
	PlayerID int `json:"player_id"`
}

func (g PlayerAcceptedRoundResponse) Type() int {
	return playerAcceptedRoundResponse
}

//////////////////////////////////////

type NewRoundResponse struct {
	ActivePlayerID int `json:"active_player_id"`
}

func (g NewRoundResponse) Type() int {
	return newRoundResponse
}

//////////////////////////////////////

type YourNewRoundResponse struct {
	Letters []game.Letter `json:"letters"`
}

func (g YourNewRoundResponse) Type() int {
	return yourNewRoundResponse
}

//////////////////////////////////////

type PlayerDeclinedWordsResponse struct {
	PlayerID   int    `json:"player_id"`
	PlayerName string `json:"player_name"`
}

func (g PlayerDeclinedWordsResponse) Type() int {
	return playerDeclinedWordsResponse
}

//////////////////////////////////////

type GameEndedResponse struct {
	PlayerPoints map[int]game.Player `json:"player_points"`
}

func (g GameEndedResponse) Type() int {
	return gameEndedResponse
}

//////////////////////////////////////

type AcceptResultedInNewRound struct{}

func (g AcceptResultedInNewRound) Type() int {
	return acceptResultedInNewRoundResponse
}

//////////////////////////////////////

type PlayerConnectionChanged struct {
	PlayerID     int  `json:"player_id"`
	Disconnected bool `json:"disconnected"`
}

func (g PlayerConnectionChanged) Type() int {
	return playerConnectionChangedResponse
}

//////////////////////////////////////

type GameStateRegenerationResponse struct {
	User                  model.User          `json:"user"`
	Players               []game.Player       `json:"players"`
	Tiles                 []game.Tile         `json:"tiles"`
	ActivePlayerID        int                 `json:"active_player_id"`
	PlayerPoints          map[int]game.Player `json:"player_points"`
	CurrentPlayerPoints   int                 `json:"current_player_points"`
	RoundFinished         bool                `json:"round_finished"`
	PlayerIDsThatAccepted []int               `json:"player_ids_that_accepted"`
	Letters               []game.Letter       `json:"letters"`
}

func (g GameStateRegenerationResponse) Type() int {
	return gameStateRegenerationResponse
}

//////////////////////////////////////

type KeepAliveResponse struct{}

func (g KeepAliveResponse) Type() int {
	return keepAliveResponse
}

//////////////////////////////////////

const (
	ServerRestarted          = 0
	ServerRestartedNameTaken = 1
	MovedToLobbyScreen       = 2
	Nothing                  = 3
)

type UserStateRegeneration struct {
	State int        `json:"state"`
	User  model.User `json:"user"`
}

func (g UserStateRegeneration) Type() int {
	return userStateRegenerationResponse
}

//////////////////////////////////////

type FinishResultedInNextRound struct{}

func (g FinishResultedInNextRound) Type() int {
	return finishResultedInNextRoundResponse
}
