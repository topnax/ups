package game

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"ups/sp/server/utils"
)

const DeskSize = 15

type Letter struct {
	Value    string `json:"value"`
	Points   int    `json:"points"`
	PlayerID int
}

type Desk struct {
	Tiles          [DeskSize][DeskSize]Tile
	Words          []Word
	CurrentLetters *LetterSet
	PlacedLetter   *LetterSet

	LetterPointsTable map[string][2]int
}

// creates a desk of default desk size
func (desk *Desk) Create() {
	letters := [DeskSize][DeskSize]Tile{}

	for row := 0; row < DeskSize; row++ {
		for column := 0; column < DeskSize; column++ {
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

	// sets the type of tiles based on the KrisKros rules
	for row := 0; row < DeskSize; row++ {
		for column := 0; column < DeskSize; column++ {
			desk.Tiles[row][column].Type = deskTypes[row][column]
		}
	}
}

// checks whether the given row and column is within desk bounds
func (desk Desk) isWithinBounds(row int, column int) bool {
	return row >= 0 && row < DeskSize && column >= 0 && column < DeskSize
}

// sets the given letter at the given row and column
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

	// call GetWordsAt so new words are found
	desk.GetWordsAt(row, column)
	return nil
}

// resets the tile at the given row and column
func (desk *Desk) ResetAt(row int, column int, playerID int) error {

	if !desk.isWithinBounds(row, column) {
		return errors.New("cannot set, out of bounds")
	}

	if desk.Tiles[row][column].Letter.PlayerID != playerID {
		return errors.New(fmt.Sprintf("Player of ID %d cannot remove player's of ID %d letters...", playerID, desk.Tiles[row][column].Letter.PlayerID))
	}

	if !desk.Tiles[row][column].Set {
		return errors.New(fmt.Sprintf("No letter set at given tile row:column %d:%d", row, column))
	}

	desk.CurrentLetters.Remove(desk.Tiles[row][column])
	desk.PlacedLetter.Remove(desk.Tiles[row][column])
	desk.Tiles[row][column].Set = false

	// all currently highlighted letters are reset
	for tile := range desk.CurrentLetters.List {
		tile.Highlighted = false
	}

	// check which words disappeared after the removal of the letter
	for tile := range desk.PlacedLetter.List {
		desk.GetWordsAt(tile.Row, tile.Column)
	}

	return nil
}

// clears the set of current letters, resets it's highlight attribute
func (desk *Desk) ClearCurrentLetters() {
	for tile, _ := range desk.CurrentLetters.List {
		desk.Tiles[tile.Row][tile.Column].Highlighted = false
	}
	desk.CurrentLetters.Clear()
}

// returns the word based on the given word meta
func (desk *Desk) GetWordAt(wordMeta WordMeta) Word {
	var content []string
	wordMultiplicand := 1

	totalPoints := 0

	// check whether the word is horizontal or vertical
	if wordMeta.RowStart == wordMeta.RowEnd {
		for column := wordMeta.ColumnStart; column < wordMeta.ColumnEnd+1; column++ {
			tile := desk.Tiles[wordMeta.RowStart][column]
			content = append(content, tile.Letter.Value)
			totalPoints += tile.getTilePoints()
			wordMultiplicand = utils.Max(wordMultiplicand, tile.getWordMultiplicand())
		}
	} else if wordMeta.ColumnStart == wordMeta.ColumnEnd {
		for row := wordMeta.RowStart; row < wordMeta.RowEnd+1; row++ {
			tile := desk.Tiles[row][wordMeta.ColumnStart]
			content = append(content, tile.Letter.Value)
			totalPoints += tile.getTilePoints()
			wordMultiplicand = utils.Max(wordMultiplicand, tile.getWordMultiplicand())
		}
	}

	totalPoints *= wordMultiplicand

	return Word{
		WordMeta: wordMeta,
		Content:  strings.Join(content[:], ""),
		Points:   totalPoints,
	}
}

// finds all word metas at the given row and column
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

	// iterate over all directions
	for index, dir := range a {
		dx := 0 + dir[0]
		dy := 0 + dir[1]

		for desk.isWithinBounds(row+dx, column+dy) && desk.Tiles[row+dx][column+dy].Set {
			desk.Tiles[row+dx][column+dy].Highlighted = true
			desk.CurrentLetters.Add(desk.Tiles[row+dx][column+dy])

			dx += dir[0]
			dy += dir[1]
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

	// if there are no other letters, highlight the only letter and mark it as a word
	if desk.isWithinBounds(row, column) && desk.Tiles[row][column].Set {
		desk.Tiles[row][column].Highlighted = true
		desk.CurrentLetters.Add(desk.Tiles[row][column])
		if mincol == maxcol && minrow == maxrow {
			words = append(words, WordMeta{
				RowStart:    row,
				ColumnStart: column,
				RowEnd:      row,
				ColumnEnd:   column,
			})
		}
	}

	return words
}

// finds all words based on the placed letters and calculates their point sum. Words are saved into a set to prevent duplicates
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
	}

	return points
}

// debug print
func (desk Desk) Print() {
	fmt.Println(" 123456789012345")
	for row := 0; row < DeskSize; row++ {
		if row < 9 {
			fmt.Print(row + 1)
		} else {
			fmt.Print(row - 9)
		}
		for column := 0; column < DeskSize; column++ {
			if desk.Tiles[row][column].Set {
				fmt.Print(strings.ToUpper(desk.Tiles[row][column].Letter.Value))

			} else {
				fmt.Print("_")
			}
		}
		fmt.Print("\n")
	}
}

// checks whether the given string is a number or not
func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
