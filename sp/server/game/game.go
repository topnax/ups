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

	PointsTable map[int]int

	PlayersThatAccepted *PlayerSet
	idInc               int
}

func (game *Game) AcceptTurn(player Player) bool {
	if player.ID == game.CurrentPlayer.ID {
		return false
	}
	game.PlayersThatAccepted.Add(player)
	return len(game.PlayersThatAccepted.List) == len(game.Players)-1
}

func (game *Game) Start() error {
	if len(game.Players) > 1 && len(game.Players) <= 4 {
		game.PointsTable = make(map[int]int)
		for _, player := range game.Players {
			game.PointsTable[player.ID] = 0
		}
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
	game.Desk.GetTotalPoints()
	game.Next()

	game.PrintPoints()

	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func (game *Game) Next() {

	game.PointsTable[game.CurrentPlayerIndex] += game.Desk.GetTotalPoints()

	if game.CurrentPlayerIndex < 0 || game.CurrentPlayerIndex == len(game.Players)-1 {
		game.CurrentPlayerIndex = 0
		game.Round++
	} else {
		game.CurrentPlayerIndex++
	}
	game.CurrentPlayer = game.Players[game.CurrentPlayerIndex]
	game.Desk.ClearCurrentWords()
	game.Desk.PlacedLetter.Clear()
	game.PlayersThatAccepted = NewPlayerSet()
}

func (game *Game) HandleSetAtEvent(event SetLetterAtEvent) error {
	return game.Desk.SetAt(event.Letter, event.Row, event.Column, event.PlayerID)
}

func (game Game) HandleResetAtEvent(event ResetAtEvent) {
	game.Desk.Tiles[event.Row][event.Column].Set = false
}

func (game Game) PrintPoints() {
	for _, player := range game.Players {
		fmt.Println(player.Name, "has", game.PointsTable[player.ID], "points")
	}
}