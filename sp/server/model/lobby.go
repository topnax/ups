package model

import "ups/sp/server/game"

type Lobby struct {
	ID      int           `json:"id"`
	Players []game.Player `json:"players"`
	Owner   game.Player   `json:"owner"`
}
