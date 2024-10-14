package main

import (
	"boredbook/explorer"
	"fmt"
	"log"
)

func main() {

	// https://gosamples.dev/read-user-input/
	fmt.Printf("What website do you want to explore? (e.g. https://google.com)\n")

	var website string
	_, err := fmt.Scanln(&website)
	if err != nil {
		log.Fatal(err)
	}

	explorer.ExploreSite(website)
}
