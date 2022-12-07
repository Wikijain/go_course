package main

import (
	"fmt"
)

func main() {
	fruitArr := []string{"apple", "orange", "grapes", "peaches"}
	fmt.Println(fruitArr)
	fruitSlice := fruitArr
	fmt.Println(fruitSlice[1:3])

}
