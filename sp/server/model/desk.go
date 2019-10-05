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
	Used  int
}

type Word struct {
	Start    [2]int
	End      [2]int
	PlayerID int
	Content  string
	Points   int
}

type Desk struct {
	Letters [deskSize][deskSize]Letter
	Words   []Word
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
				Used:  0,
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

	var content []string

	if rowStart == rowEnd {
		for column := columnStart; column < columnEnd+1; column++ {
			if desk.Letters[rowStart][column].Used >= 2 {
				return errors.New("a word cannot be used more than two times ")
			} else if !desk.Letters[rowStart][column].Set {
				return errors.New("a letter has to be set everywhere inside the range")
			}
		}

		if desk.Letters[rowStart][columnStart-1].Set || desk.Letters[rowStart][columnEnd+1].Set {
			return errors.New("no letter can't be after or before the word")
		}

		for column := columnStart; column < columnEnd+1; column++ {
			content = append(content, desk.Letters[rowStart][column].Value)
			desk.Letters[rowStart][column].Used++
		}
	}

	if columnStart == columnEnd {
		for row := rowStart; row < rowEnd+1; row++ {
			if desk.Letters[row][columnStart].Used >= 2 {
				return errors.New("a word cannot be used more than two times ")
			} else if !desk.Letters[row][columnStart].Set {
				return errors.New("a letter has to be set everywhere inside the range")
			}
		}

		if desk.Letters[rowEnd+1][columnStart].Set || desk.Letters[rowStart-1][columnStart].Set {
			return errors.New("no letter can't be after or before the word")
		}

		for row := rowStart; row < rowEnd+1; row++ {
			content = append(content, desk.Letters[row][columnStart].Value)
			desk.Letters[row][columnStart].Used++
		}
	}

	desk.Words = append(desk.Words, Word{
		Start:    [2]int{rowStart, columnStart},
		End:      [2]int{rowEnd, columnEnd},
		PlayerID: playerId,
		Content:  strings.Join(content[:], ""),
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
