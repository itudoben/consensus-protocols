package main

import (
	"fmt"
	"log"

	// go mod edit -replace example.com/greetings=../greetings
	// It allows to use the local code when developing and before publishing it.
	"itudoben.io/greetings"
	"rsc.io/quote"
)

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(log.Ldate | log.Ltime)

	log.Print("testing the multiple return values function")
	message, err := greetings.Hello("George")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
	fmt.Println("Hello, World!")
	fmt.Println(quote.Go())
}
