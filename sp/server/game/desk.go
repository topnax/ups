package game

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"ups/sp/server/utils"
)

const deskSize = 15

type Letter struct {
	Value    string `json:"value"`
	Points   int    `json:"points"`
	PlayerID int
}

type Desk struct {
	Tiles          [deskSize][deskSize]Tile
	Words          []Word
	CurrentLetters *LetterSet
	PlacedLetter   *LetterSet

	LetterPointsTable map[string][2]int
}

func (desk *Desk) Create() {
	letters := [deskSize][deskSize]Tile{}

	for row := 0; row < deskSize; row++ {
		for column := 0; column < deskSize; column++ {
			letters[row][column] = Tile{
				Set:    false,
				Type:   0,
				Row:    row,
				Column: column,
			}
		}
	}

	desk.Tiles = letters
	desk.CurrentLetters = NewSet()
	desk.PlacedLetter = NewSet()
	desk.LetterPointsTable = GetLetterPointsTable()

	deskTypes := GetDeskTileTypes()

	for row := 0; row < deskSize; row++ {
		for column := 0; column < deskSize; column++ {
			desk.Tiles[row][column].Type = deskTypes[row][column]
		}
	}
}

func (desk Desk) isWithinBounds(row int, column int) bool {
	return row >= 0 && row < deskSize && column >= 0 && column < deskSize
}

func (desk *Desk) SetAt(letter string, row int, column int, playerID int) error {
	letterPoints, exists := desk.LetterPointsTable[strings.ToLower(letter)]
	if !exists {
		return errors.New("Letter " + letter + "not found in the letter table")
	}

	if isNumber(letter) {
		return errors.New("Only letters are allowed")
	}
	if !desk.isWithinBounds(row, column) {
		return errors.New("cannot set, out of bounds")
	}

	if desk.Tiles[row][column].Set {
		return errors.New("letter already set")
	}

	desk.Tiles[row][column].Letter = Letter{
		Points:   letterPoints[0],
		Value:    letter,
		PlayerID: playerID,
	}
	desk.Tiles[row][column].Set = true
	desk.Tiles[row][column].Highlighted = true
	desk.CurrentLetters.Add(desk.Tiles[row][column])
	desk.PlacedLetter.Add(desk.Tiles[row][column])

	desk.GetWordsAt(row, column)
	return nil
}

func (desk *Desk) ResetAt(row int, column int, playerID int) error {

	if !desk.isWithinBounds(row, column) {
		return errors.New("cannot set, out of bounds")
	}

	if desk.Tiles[row][column].Letter.PlayerID != playerID {
		return errors.New(fmt.Sprintf("Player #%d cannot remove player's #%d letters...", playerID, desk.Tiles[row][column].Letter.PlayerID))
	}

	if !desk.Tiles[row][column].Set {
		return errors.New(fmt.Sprintf("No letter set at given tile row:column %d:%d", row, column))
	}

	desk.Tiles[row][column].Set = false

	desk.CurrentLetters.Remove(desk.Tiles[row][column])
	desk.PlacedLetter.Remove(desk.Tiles[row][column])
	return nil
}

func (desk *Desk) ClearCurrentWords() {
	desk.CurrentLetters.Clear()
}

func (desk *Desk) GetWordAt(wordMeta WordMeta) Word {
	var content []string
	wordMultiplicand := 1

	totalPoints := 0

	if wordMeta.RowStart == wordMeta.RowEnd {
		for column := wordMeta.ColumnStart; column < wordMeta.ColumnEnd+1; column++ {
			tile := desk.Tiles[wordMeta.RowStart][column]
			content = append(content, tile.Letter.Value)
			totalPoints += tile.getTilePoints()
			wordMultiplicand = utils.Max(wordMultiplicand, tile.getWordMultiplicand())
		}
	}

	if wordMeta.ColumnStart == wordMeta.ColumnEnd {
		for row := wordMeta.RowStart; row < wordMeta.RowEnd+1; row++ {
			tile := desk.Tiles[row][wordMeta.ColumnStart]
			content = append(content, tile.Letter.Value)
			totalPoints += tile.getTilePoints()
			wordMultiplicand = utils.Max(wordMultiplicand, tile.getWordMultiplicand())
		}
	}

	fmt.Println(strings.Join(content[:], ""), wordMultiplicand, totalPoints)
	totalPoints *= wordMultiplicand

	return Word{
		WordMeta: wordMeta,
		Content:  strings.Join(content[:], ""),
		Points:   totalPoints,
	}
}

func (desk *Desk) GetWordsAt(row int, column int) []WordMeta {

	// directions
	a := [][]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}

	mincol := -1
	maxcol := -1

	minrow := -1
	maxrow := -1

	for index, c := range a {
		dx := 0 + c[0]
		dy := 0 + c[1]

		for desk.isWithinBounds(row+dx, column+dy) && desk.Tiles[row+dx][column+dy].Set {
			desk.CurrentLetters.Add(desk.Tiles[row+dx][column+dy])
			desk.Tiles[row+dx][column+dy].Highlighted = true

			dx += c[0]
			dy += c[1]
		}

		switch index {
		case 0:
			maxcol = column + dy - 1
			break
		case 1:
			mincol = column + dy + 1
			break
		case 2:
			maxrow = row + dx - 1
			break
		case 3:
			minrow = row + dx + 1
			break
		}
	}
	var words []WordMeta

	if mincol != maxcol {
		words = append(words, WordMeta{
			RowStart:    row,
			ColumnStart: mincol,
			RowEnd:      row,
			ColumnEnd:   maxcol,
		})
	}

	if minrow != maxrow {
		words = append(words, WordMeta{
			RowStart:    minrow,
			ColumnStart: column,
			RowEnd:      maxrow,
			ColumnEnd:   column,
		})
	}

	return words
}

func (desk Desk) GetTotalPoints() int {
	wordMap := NewWordMetaSet()

	var words []Word

	for letter, _ := range desk.PlacedLetter.List {
		for _, wordMeta := range desk.GetWordsAt(letter.Row, letter.Column) {
			if !wordMap.Has(wordMeta) {
				words = append(words, desk.GetWordAt(wordMeta))
				wordMap.Add(wordMeta)
			}
		}
	}

	points := 0

	for _, word := range words {
		points += word.Points
		fmt.Println("Word:", word.Points, word.Content, "at", word.WordMeta.RowStart, word.WordMeta.ColumnStart, "-", word.WordMeta.RowEnd, word.WordMeta.ColumnEnd)
	}

	fmt.Println("Total points:", points)
	return points

}

func (desk Desk) Print() {
	fmt.Println(" 123456789012345")
	for row := 0; row < deskSize; row++ {
		if row < 9 {
			fmt.Print(row + 1)
		} else {
			fmt.Print(row - 9)
		}
		for column := 0; column < deskSize; column++ {
			if desk.Tiles[row][column].Set {
				fmt.Print(strings.ToUpper(desk.Tiles[row][column].Letter.Value))

			} else {
				fmt.Print("_")
			}
		}
		fmt.Print("\n")
	}
}

func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
