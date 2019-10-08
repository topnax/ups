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

type Letter struct {
	Row           int
	Column        int
	Set           bool
	Type          int
	Value         string
	Used          int
	CurrentlyUsed bool
}

type WordMeta struct {
	RowStart    int
	ColumnStart int
	RowEnd      int
	ColumnEnd   int
}

type Word struct {
	Start    [2]int
	End      [2]int
	PlayerID int
	Content  string
	Points   int
}

type Desk struct {
	Letters        [deskSize][deskSize]Letter
	Words          []Word
	CurrentLetters *LetterSet
	PlacedLetter   *LetterSet

	lastRow int
	lastCol int
}

type KrisKrosDesk interface {
	Create()
	Print() func()
	SetAt(letter Letter, x, y int) bool
}

func (desk *Desk) Create() {
	letters := [deskSize][deskSize]Letter{}

	for row := 0; row < deskSize; row++ {
		for column := 0; column < deskSize; column++ {
			letters[row][column] = Letter{
				Set:    false,
				Type:   0,
				Used:   0,
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

	a := [][]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}

	fmt.Printf("\nSet: %v:%v - %s\n", row, column, letter)
	for _, c := range a {
		dx := 0 + c[0]
		dy := 0 + c[1]
		for desk.isWithinBounds(row+dx, column+dy) && desk.Letters[row+dx][column+dy].Set {
			desk.CurrentLetters.Add(desk.Letters[row+dx][column+dy])
			dx += c[0]
			dy += c[1]
		}
	}

	//fmt.Println("Used letters:")
	//for _, value := range usedLetters {
	//	fmt.Println(value)
	//}

	desk.Letters[row][column].Value = letter
	desk.Letters[row][column].Set = true
	fmt.Println("setting", desk.Letters[row][column].Value, "at", row, column)
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

func (desk *Desk) GetWordAt(rowStart int, columnStart int, rowEnd int, columnEnd int) string {
	var content []string

	if rowStart == rowEnd {
		for column := columnStart; column < columnEnd+1; column++ {
			content = append(content, desk.Letters[rowStart][column].Value)
		}
	}

	if columnStart == columnEnd {
		for row := rowStart; row < rowEnd+1; row++ {
			content = append(content, desk.Letters[row][columnStart].Value)
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

func (desk *Desk) GetWordsAt(row int, column int) [][]int {
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
	var words [][]int

	words = append(words, []int{mincol, maxcol, row})

	words = append(words, []int{minrow, maxrow, column})

	println()
	return words
}

func (desk Desk) GetTotalPoints() {
	wordMap := NewWordMetaSet()

	var words []Word

	for letter, _ := range desk.PlacedLetter.List {
		fmt.Println(letter.Value)
		for i, positions := range desk.GetWordsAt(letter.Row, letter.Column) {
			if positions[0] != -1 && positions[1] != -1 && positions[2] != -1 && positions[0] != positions[1] {
				if i == 0 {
					wordMeta := WordMeta{
						RowStart:    positions[2],
						ColumnStart: positions[0],
						RowEnd:      positions[2],
						ColumnEnd:   positions[1],
					}
					if !wordMap.Has(wordMeta) {
						word := desk.GetWordAt(positions[2], positions[0], positions[2], positions[1])
						words = append(words, Word{
							Start:    [2]int{positions[1], positions[0]},
							End:      [2]int{positions[2], positions[0]},
							PlayerID: 0,
							Content:  word,
							Points:   0,
						})
						wordMap.Add(wordMeta)
					}
				}
				if i == 1 {
					wordMeta := WordMeta{
						RowStart:    positions[0],
						ColumnStart: positions[2],
						RowEnd:      positions[1],
						ColumnEnd:   positions[2],
					}
					if !wordMap.Has(wordMeta) {
						word := desk.GetWordAt(positions[0], positions[2], positions[1], positions[2])
						words = append(words, Word{
							Start:    [2]int{positions[0], positions[1]},
							End:      [2]int{positions[0], positions[2]},
							PlayerID: 0,
							Content:  word,
							Points:   0,
						})
						wordMap.Add(wordMeta)
					}

				}
			}

		}
	}

	for _, word := range words {
		fmt.Println("Word:", word.Content, "at", word.Start[0], word.Start[1], "-", word.End[0], word.End[1])
	}
	//multiplicators := map[Letter]int
	//totalPoints := 0
	//for letter, _ := range desk.CurrentLetters.List {
	//	a := [][] int{
	//		{0, 0},
	//		{0, 1},
	//		{0, -1},
	//		{1, 0},
	//		{-1, 0},
	//	}
	//	row := letter.Row
	//	column := letter.Column
	//
	//	fmt.Printf("\nSet: %v:%v - %s\n", row, column, letter)
	//	for _, c := range a {
	//		dx := 0 + c[0]
	//		dy := 0 + c[1]
	//		for desk.isWithinBounds(row+dx, column+dy) && desk.Letters[row+dx][column+dy].Set {
	//			desk.CurrentLetters.Add(desk.Letters[row+dx][column+dy])
	//			dx += c[0]
	//			dy += c[1]
	//		}
	//	}
	//
	//	switch letter.Type {
	//
	//	}
	//}
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
				if desk.Letters[row][column].CurrentlyUsed {
					fmt.Print(strings.ToUpper(desk.Letters[row][column].Value))
				} else {
					fmt.Print(strings.ToLower(desk.Letters[row][column].Value))
				}
			} else {
				if desk.Letters[row][column].CurrentlyUsed {
					fmt.Print("&")
				} else {
					fmt.Print("_")
				}
			}
		}
		fmt.Print("\n")
	}

	if len(desk.Words) > 0 {
		fmt.Println("Words:")
		for _, word := range desk.Words {
			fmt.Println(word.Content, "by user", word.PlayerID, "at", word.Start[0], word.Start[1], "to", word.End[0], word.End[1])
		}
	}

}

func GetDesk() Desk {
	desk := Desk{}
	desk.Create()
	return desk
}
