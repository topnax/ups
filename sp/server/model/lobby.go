package model

import "ups/sp/server/game"

type Lobby struct {
	ID      int           `json:"id"`
	Players []game.Player `json:"players"`
	Owner   game.Player   `json:"owner"`
}

// a method that returns true when start of the lobby is possible
func (lobby Lobby) IsStartPossible() bool {
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
