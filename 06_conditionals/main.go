package main

import "fmt"

func main() {
	x := 10
	y := 20

	if x < y {
		fmt.Println("x < y")
	} else if x == y {
		fmt.Println("x == y")
	} else {
		fmt.Println("x > y")
	}
}
