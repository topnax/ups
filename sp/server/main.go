package main

import (
<<<<<<< HEAD
	"fmt"
=======
	"ups/sp/server/game"
>>>>>>> f868af92b9420e7f66281cabe0a1216a9f0d0009
)

func main() {

<<<<<<< HEAD
	var age int

	fmt.Println(age)

	fmt.Println("My name is", "Foo", "and my age is of", 10, "years and the current hour is", 19)

	name := "foo"

	name = "bar"

	fmt.Println(name)
	fmt.Println(name)
	fmt.Println(name)

	desk.

	area, perimeter := rectProps(5, 6)

	fmt.Println("The area is", area, "and the perimeter is", perimeter)
	printPyra(20)
	arr := [...]string {"Tom","Pavel", "Dan", "Petr", "Trtek"}

	for index, element := range arr[1::4]{
		fmt.Println(index, element)
	}
}

func rectProps(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = (length + width) * 2
	return //no explicit return value
}

func printPyra(height int) {
	fmt.Println("Printing pyramid of height", height)
	floor := height*2 - 1
	for i := 0; i < height; i++ {
		air := height - 1 - i
		for h := 0; h < floor; h++ {
			if h < air || h > floor-air-1 {
				fmt.Print(" ")
			} else {
				if h % 2 == 0 {
					fmt.Print("O")
				} else {
					fmt.Print("0")
				}
			}
		}
		fmt.Print("\n")
	}
=======
	currentGame := game.Game{}
	currentGame.AddPlayer("Pavel")
	currentGame.AddPlayer("Tomáš")
	currentGame.AddPlayer("Fanda")

	currentGame.Start()

	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      1,
		Column:   1,
		Letter:   "S",
	})
	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      2,
		Column:   1,
		Letter:   "E",
	})
	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      3,
		Column:   1,
		Letter:   "X",
	})
	currentGame.Print()

	currentGame.Next()
	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      1,
		Column:   3,
		Letter:   "P",
	})
	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      2,
		Column:   3,
		Letter:   "E",
	})
	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      3,
		Column:   3,
		Letter:   "S",
	})

	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      2,
		Column:   2,
		Letter:   "Y",
	})
	currentGame.Print()

	currentGame.Next()

	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      2,
		Column:   4,
		Letter:   "S",
	})

	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      4,
		Column:   1,
		Letter:   "Y",
	})

	_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: currentGame.CurrentPlayer.ID,
		Row:      5,
		Column:   1,
		Letter:   "H",
	})

	currentGame.Print()

	currentGame.Desk.GetTotalPoints()

	////
	//desk := model.GetDesk()
	//
	//desk.Print()
	//desk.SetAt("A", 0, 1)
	//desk.SetAt("B", 2, 1)
	//desk.SetAt("B", 2, 1)
	//desk.SetAt("CH", 14, 14)
	//fmt.Println("Out of bounds")
	//desk.SetAt("K", 3, 3)
	//desk.SetAt("R", 3, 4)
	//desk.SetAt("I", 3, 5)
	//desk.SetAt("S", 3, 6)
	//
	//desk.SetAt("K", 2, 4)
	//
	//desk.SetAt("O", 4, 4)
	//desk.SetAt("S", 5, 4)
	//
	//desk.SetAt("Z", 14, 15)
	//desk.SetAt("Z", 15, 14)
	//desk.SetAt("CH", 14, 14)
	//desk.SetAt("131", 13, 13)
	//desk.SetAt("1", 12, 13)
	//desk.SetAt("Z", -1, 15)
	//desk.SetAt("Z", 15, -1)
	//desk.SetAt("Z", 15, -1)
	//
	//error := desk.SetWordAt(3, 3, 3, 6, 1)
	//error2 := desk.SetWordAt(2, 4, 5, 4, 1)
	//
	////desk.SetAt("B", 16,15)
	////desk.SetAt("B", 15,15)
	//desk.Print()
	//
	//if error != nil {
	//	fmt.Println(error)
	//}
	//
	//if error2 != nil {
	//	fmt.Println(error2)
	//}
>>>>>>> f868af92b9420e7f66281cabe0a1216a9f0d0009
}
