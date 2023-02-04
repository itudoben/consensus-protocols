package main

import (
	"fmt"

	"itudoben.io/greetings"
	"rsc.io/quote"
)

func main() {
	message := greetings.Hello("Gladys")
	fmt.Println(message)
	fmt.Println("Hello, World!")
	fmt.Println(quote.Go())
}
