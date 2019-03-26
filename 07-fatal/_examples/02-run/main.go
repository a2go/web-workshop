package main

import (
	"errors"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("Start program")
	defer log.Println("End program")

	if err := open(); err != nil {
		return err
	}

	return nil
}

func open() error {
	return errors.New("whoops")
}
