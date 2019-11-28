package model

import "ups/sp/server/game"

type Lobby struct {
	ID      int           `json:"id"`
	Players []game.Player `json:"players"`
	Owner   game.Player   `json:"owner"`
}

func (lobby Lobby) IsStartable() bool {
	playerCount := len(lobby.Players)
	if playerCount > 1 {
		readyPlayers := 0
		for _, player := range lobby.Players {
			if player.Ready {
				readyPlayers++
			}
		}

		if readyPlayers == playerCount {
			return true
		}
	}
	return false
}
