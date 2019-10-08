package model

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const deskSize = 15

const ( // iota is reset to 0
	BASIC             = 0
	MULTIPLY_WORD_2   = 1
	MULTIPLY_WORD_3   = 2
	MULTIPLY_LETTER_2 = 3
	MULTIPLY_LETTER_3 = 4
)

type Tile struct {
	Row    int
	Column int
	Set    bool
	Type   int
	Value  string
}

type WordMeta struct {
	RowStart    int
	ColumnStart int
	RowEnd      int
	ColumnEnd   int
}

type Word struct {
	WordMeta WordMeta
	PlayerID int
	Content  string
	Points   int
}

type Desk struct {
	Letters        [deskSize][deskSize]Tile
	Words          []Word
	CurrentLetters *LetterSet
	PlacedLetter   *LetterSet

	lastRow int
	lastCol int
}

type KrisKrosDesk interface {
	Create()
	Print() func()
	SetAt(letter Tile, x, y int) bool
}

func (desk *Desk) Create() {
	letters := [deskSize][deskSize]Tile{}

	for row := 0; row < deskSize; row++ {
		for column := 0; column < deskSize; column++ {
			letters[row][column] = Tile{
				Set:    false,
				Type:   0,
				Value:  "",
				Row:    row,
				Column: column,
			}
		}
	}
	desk.Letters = letters
	desk.CurrentLetters = NewSet()
	desk.PlacedLetter = NewSet()
}

func (desk Desk) isWithinBounds(row int, column int) bool {
	return row >= 0 && row < deskSize && column >= 0 && column < deskSize
}

func (desk *Desk) SetAt(letter string, row int, column int) error {
	if (letter != "CH" && len(letter) > 1) || isNumber(letter) {
		return errors.New("cannot set a letter longer than 1. 'CH' is an exception. Only letters are allowed")
	}
	if !desk.isWithinBounds(row, column) {
		return errors.New("cannot set, out of bounds")
	}

	if desk.Letters[row][column].Set {
		return errors.New("letter already set")
	}

	desk.Letters[row][column].Value = letter
	desk.Letters[row][column].Set = true
	desk.CurrentLetters.Add(desk.Letters[row][column])
	desk.PlacedLetter.Add(desk.Letters[row][column])
	return nil
}

func (desk *Desk) ClearCurrentWords() {
	desk.CurrentLetters.Clear()
}

func (desk *Desk) ClearLast() {
	desk.lastCol = -1
	desk.lastRow = -1
}

func (desk *Desk) GetWordAt(wordMeta WordMeta) string {
	var content []string

	if wordMeta.RowStart == wordMeta.RowEnd {
		for column := wordMeta.ColumnStart; column < wordMeta.ColumnEnd+1; column++ {
			content = append(content, desk.Letters[wordMeta.RowStart][column].Value)
		}
	}

	if wordMeta.ColumnStart == wordMeta.ColumnEnd {
		for row := wordMeta.RowStart; row < wordMeta.RowEnd+1; row++ {
			content = append(content, desk.Letters[row][wordMeta.ColumnStart].Value)
		}
	}
	return strings.Join(content[:], "")
}

func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
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

		for desk.isWithinBounds(row+dx, column+dy) && desk.Letters[row+dx][column+dy].Set {
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

func (desk Desk) GetTotalPoints() {
	wordMap := NewWordMetaSet()

	var words []Word

	for letter, _ := range desk.PlacedLetter.List {
		for _, wordMeta := range desk.GetWordsAt(letter.Row, letter.Column) {
			if !wordMap.Has(wordMeta) {
				word := desk.GetWordAt(wordMeta)
				words = append(words, Word{
					WordMeta: wordMeta,
					PlayerID: 0,
					Content:  word,
					Points:   0,
				})
				wordMap.Add(wordMeta)
			}

		}
	}

	for _, word := range words {
		fmt.Println("Word:", word.Content, "at", word.WordMeta.RowStart, word.WordMeta.ColumnStart, "-", word.WordMeta.RowEnd, word.WordMeta.ColumnEnd)
	}

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
			if desk.Letters[row][column].Set {
				fmt.Print(strings.ToUpper(desk.Letters[row][column].Value))

			} else {
				fmt.Print("_")
			}
		}
		fmt.Print("\n")
	}

	if len(desk.Words) > 0 {
		//fmt.Println("Words:")
		//for _, word := range desk.Words {
		//	fmt.Println(word.Content, "by user", word.PlayerID, "at", word.Start[0], word.Start[1], "to", word.End[0], word.End[1])
		//}
	}

}

func GetDesk() Desk {
	desk := Desk{}
	desk.Create()
	return desk
}

//a := [][]int{
//{0, 1},
//{0, -1},
//{1, 0},
//{-1, 0},
//}
//
//for _, c := range a {
//dx := 0 + c[0]
//dy := 0 + c[1]
//for desk.isWithinBounds(row+dx, column+dy) && desk.Letters[row+dx][column+dy].Set {
//desk.CurrentLetters.Add(desk.Letters[row+dx][column+dy])
//dx += c[0]
//dy += c[1]
//}
//}
