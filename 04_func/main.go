package main

import "fmt"

func greeting(s string) string {
	return "Hello " + s
}
func fibbo(n int) int {
	if n < 2 {
		return n
	}
	return fibbo(n-1) + fibbo(n-2)
}
func main() {
	fmt.Println(greeting("Harsh"))
	fmt.Println(fibbo(0))
}
