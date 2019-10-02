package main

import (
	"fmt"
	"ups/sp/server/model"
)


func main() {

	//
	desk := model.GetDesk()

	desk.Print()
	desk.SetAt("A", 0,1)
	desk.SetAt("B", 2,1)
	desk.SetAt("B", 2,1)
	desk.SetAt("CH", 14,14)
	fmt.Println("Out of bounds")
	desk.SetAt("Z", 15,15)
	desk.SetAt("Z", 14,15)
	desk.SetAt("Z", 15,14)
	desk.SetAt("CH", 14,14)
	desk.SetAt("131", 13,13)
	desk.SetAt("1", 12,13)
	desk.SetAt("Z", -1,15)
	desk.SetAt("Z", 15,-1)
	desk.SetAt("Z", 15,-1)
	//desk.SetAt("B", 16,15)
	//desk.SetAt("B", 15,15)
	desk.Print()
}
