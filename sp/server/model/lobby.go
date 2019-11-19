package model

import "ups/sp/server/game"

type Lobby struct {
	ID      int
	Players []game.Player
	Owner   game.Player
}
