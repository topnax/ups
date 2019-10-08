package game

import (
	"errors"
	"fmt"
)

type Game struct {
	Desk               Desk
	Players            []Player
	CurrentPlayer      Player
	CurrentPlayerIndex int
	Round              int

	idInc int
}

func (game *Game) Start() error {
	if len(game.Players) > 1 && len(game.Players) <= 4 {
		game.Desk.Create()
		game.CurrentPlayerIndex = -1
		game.Next()
		game.Round = 1
		return nil
	} else {
		return errors.New("a game should consist of 2 to 4 players")
	}
}

func (game *Game) AddPlayer(name string) {
	if len(game.Players) < 4 {
		game.Players = append(game.Players, Player{
			Name: name,
			ID:   game.idInc,
		})
		game.idInc += 1
	}
}

func (game *Game) Print() {
	fmt.Println("Game status:")
	fmt.Println("Round:", game.Round)
	fmt.Println("Players:")
	for _, player := range game.Players {
		if player == game.CurrentPlayer {
			fmt.Printf("> #%v  %s\n", player.ID, player.Name)
		} else {
			fmt.Printf("#%v  %s\n", player.ID, player.Name)
		}

	}

	game.Desk.Print()

	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func (game *Game) Next() {
	if game.CurrentPlayerIndex < 0 || game.CurrentPlayerIndex == len(game.Players)-1 {
		game.CurrentPlayerIndex = 0
		game.Round++
	} else {
		game.CurrentPlayerIndex++
	}
	game.CurrentPlayer = game.Players[game.CurrentPlayerIndex]
	game.Desk.ClearCurrentWords()
	game.Desk.PlacedLetter.Clear()
}

func (game *Game) HandleSetAtEvent(event SetLetterAtEvent) error {
	return game.Desk.SetAt(event.Letter, event.Row, event.Column, event.PlayerID)
}

func (game Game) HandleResetAtEvent(event ResetAtEvent) {
	game.Desk.Tiles[event.Row][event.Column].Set = false
}
