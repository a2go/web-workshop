package main

import (
	"errors"
	"log"
)

func main() {
	log.Println("Start program")
	defer log.Println("End program")

	if err := open(); err != nil {
		log.Fatal(err)
	}
}

func open() error {
	return errors.New("whoops")
}
