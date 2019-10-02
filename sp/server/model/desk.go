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
	Set   bool
	Type  int
	Value string
	Used int
}

type Word struct {
	Start [2] int
	End [2] int
	PlayerID int
}

type Desk struct {
	Letters [deskSize][deskSize] Letter
	Words [] Word
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
				Set:   false,
				Type:  0,
				Used: 0,
				Value: "",
			}
		}
	}
	desk.Letters = letters
}

func (desk *Desk) SetAt(letter string, row int, column int) error {
	if (letter != "CH" && len(letter) > 1) || isNumber(letter) {
		return errors.New("cannot set a letter longer than 1. 'CH' is an exception. Only letters are allowed")
	}
	if row >= deskSize || column >= deskSize || row < 0 || column < 0 {
		return errors.New("cannot set, out of bounds")
	}
	if !desk.Letters[row][column].Set {
		desk.Letters[row][column].Value = letter
		desk.Letters[row][column].Set = true
		return nil
	}
	return errors.New("letter already set")
}

func (desk *Desk) SetWordAt(rowStart int, columnStart int, rowEnd int, columnEnd int, playerId int) error {

	if rowStart == rowEnd && columnStart == columnEnd {
		return errors.New("word has to have at least one letter")
	}

	if rowStart != rowEnd && columnStart != columnEnd {
		return errors.New("a word mustn't be diagonal")
	}

	if rowStart == rowEnd {
		for column := columnStart; column < columnEnd; column++ {
			if desk.Letters[rowStart][column].Used >= 2 {
				return errors.New("a word cannot be used more than two times ")
			}
		}

		for column := columnStart; column < columnEnd; column++ {
			desk.Letters[rowStart][column].Used++
		}
	}

	if columnStart == columnEnd {
		for row := rowStart; row < rowEnd; row++ {
			if desk.Letters[row][columnStart].Used >= 2 {
				return errors.New("a word cannot be used more than two times ")
			}
		}

		for row := rowStart; row < rowEnd; row++ {
			desk.Letters[row][columnStart].Used++
		}
	}

	desk.Words = append(desk.Words, Word{
		Start:    [2]int{rowStart, columnStart},
		End:      [2]int{rowEnd, columnStart},
		PlayerID: playerId,
	})

	return nil
}

func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func (desk Desk) Print() {
	for row := 0; row < deskSize; row++ {
		for column := 0; column < deskSize; column++ {
			if desk.Letters[row][column].Set {
				fmt.Print(strings.ToUpper(desk.Letters[row][column].Value))
			} else {
				fmt.Print("_")
			}
		}
		fmt.Print("\n")

	}
}

func GetDesk() Desk {
	desk := Desk{}
	desk.Create()
	return desk
}
