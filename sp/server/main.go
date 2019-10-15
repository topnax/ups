package main

import (
	"fmt"
)

func main() {

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
}
