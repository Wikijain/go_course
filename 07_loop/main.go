package main

import "fmt"

func main() {
	//long method loop
	i := 1
	for i <= 10 {
		fmt.Println(i)
		i++
	}

	//short method loop

	for i := 1; i <= 10; i++ {
		fmt.Printf("Number %d\n", i)
	}

	//Fizzbuzz
	for i := 1; i < 20; i++ {
		if i%15 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}
