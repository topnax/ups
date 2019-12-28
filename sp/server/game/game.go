package game

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Game struct {
	Desk               Desk
	Players            []Player
	PlayersMap         map[int]Player
	CurrentPlayer      Player
	CurrentPlayerIndex int
	Round              int

	PlayerIdToPlayerBag map[int][]Letter
	letterPointsTable   map[string][2]int

	PointsTable map[int]int

	PlayersThatAccepted *PlayerSet
	idInc               int

	EmptyRounds int

	letterBag []string
}

func (game Game) ActivePlayerCount() int {
	count := 0
	for _, player := range game.PlayersMap {
		if !player.Disconnected {
			count++
		}
	}
	return count
}

func (game *Game) AcceptTurn(player Player) bool {
	if player.ID == game.CurrentPlayer.ID {
		return false
	}
	game.PlayersThatAccepted.Add(player)

	numberOfPlayersThatMustHaveAccepted := 0

	for _, player := range game.Players {
		if !player.Disconnected && player.ID != game.CurrentPlayer.ID {
			numberOfPlayersThatMustHaveAccepted++
		}
	}

	return len(game.PlayersThatAccepted.List) >= numberOfPlayersThatMustHaveAccepted
}

func getLettersFromBag(bag []string, requested int, letterPointsTable map[string][2]int) ([]Letter, []string) {
	var randomLetters []Letter
	rand.Seed(time.Now().UnixNano())
	if len(bag) >= requested {
		for i := 0; i < requested; i++ {
			index := rand.Intn(len(bag))
			randomLetters = append(randomLetters, Letter{
				Value:    bag[index],
				Points:   letterPointsTable[bag[index]][0],
				PlayerID: 0,
			})
			fmt.Println("From index", index, "letter:", bag[index])
			bag = append(bag[:i], bag[i+1:]...)
		}
	}
	logrus.Debugf("Taking %d letters out of the bag, %d left in the great bag!", requested, len(bag))
	letters := ""
	for _, letter := range randomLetters {
		letters += " " + letter.Value
	}
	logrus.Debugln("Took out of bag: š", letters)
	return randomLetters, bag
}

func (game *Game) Start() error {
	if len(game.Players) > 1 && len(game.Players) <= 4 {
		game.PointsTable = make(map[int]int)
		game.PlayerIdToPlayerBag = make(map[int][]Letter)
		game.PlayersMap = make(map[int]Player)
		game.letterPointsTable = GetLetterPointsTable()
		game.letterBag = generateLetterBag()
		for _, player := range game.Players {
			game.PlayersMap[player.ID] = player
			game.PointsTable[player.ID] = 0
			letters, bag := getLettersFromBag(game.letterBag, 8, game.letterPointsTable)
			game.letterBag = bag
			game.PlayerIdToPlayerBag[player.ID] = letters
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

func generateLetterBag() []string {
	var letterBag []string
	for letter, info := range GetLetterPointsTable() {
		for i := 0; i < info[1]; i++ {
			letterBag = append(letterBag, letter)
		}
	}
	return letterBag
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
	//for _, player := range game.Players {
	//if player == game.CurrentPlayer {
	//	fmt.Printf("> #%v  %s\n", player.ID, player.Name)
	//} else {
	//	fmt.Printf("#%v  %s\n", player.ID, player.Name)
	//}

	//}

	game.Desk.Print()
	game.Desk.GetTotalPoints()
	//game.Next()
	//
	//game.PrintPoints()

	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func (game *Game) Next() {

	game.PointsTable[game.CurrentPlayer.ID] += game.Desk.GetTotalPoints()

	// TODO what to do when all players leave.

	for {
		if game.CurrentPlayerIndex < 0 || game.CurrentPlayerIndex == len(game.Players)-1 {
			game.CurrentPlayerIndex = 0
			game.Round++
		} else {
			game.CurrentPlayerIndex++
		}
		game.CurrentPlayer = game.Players[game.CurrentPlayerIndex]
		if !game.CurrentPlayer.Disconnected {
			break
		}

	}

	for tile, _ := range game.Desk.CurrentLetters.List {
		tile.Highlighted = false
	}

	game.Desk.ClearCurrentWords()
	game.Desk.PlacedLetter.Clear()
	game.PlayersThatAccepted = NewPlayerSet()

	newLetters, bag := getLettersFromBag(game.letterBag, 8-len(game.PlayerIdToPlayerBag[game.CurrentPlayer.ID]), game.letterPointsTable)

	game.letterBag = bag
	game.PlayerIdToPlayerBag[game.CurrentPlayer.ID] = append(game.PlayerIdToPlayerBag[game.CurrentPlayer.ID], newLetters...)
}

func (game *Game) HandleSetAtEvent(event SetLetterAtEvent) error {
	if event.PlayerID == game.CurrentPlayer.ID {
		for index, letter := range game.PlayerIdToPlayerBag[event.PlayerID] {
			if letter.Value == event.Letter {
				game.PlayerIdToPlayerBag[event.PlayerID] = append(game.PlayerIdToPlayerBag[event.PlayerID][:index], game.PlayerIdToPlayerBag[event.PlayerID][index+1:]...)
				return game.Desk.SetAt(event.Letter, event.Row, event.Column, event.PlayerID)
			}
		}
		return errors.New(fmt.Sprintf("Player does not have letter '%s' in his bag...", event.Letter))
	}
	return errors.New(fmt.Sprintf("It's not players of ID %d turn", event.PlayerID))
}

func (game *Game) HandleResetAtEvent(event ResetAtEvent) error {
	err := game.Desk.ResetAt(event.Row, event.Column, event.PlayerID)

	if err == nil {
		game.PlayerIdToPlayerBag[event.PlayerID] = append(game.PlayerIdToPlayerBag[event.PlayerID], game.Desk.Tiles[event.Row][event.Column].Letter)
	}

	return err
}

func (game Game) PrintPoints() {
	for _, player := range game.Players {
		fmt.Println(player.Name, "has", game.PointsTable[player.ID], "points")
	}
}

func (game *Game) WordsDeclined() {
	game.PlayersThatAccepted = NewPlayerSet()
}
