package main

import (
	"log"
	"time"
)

func main() {
	x := 0

	go func() {
		y0 := x
		y0++

		log.Printf("y0: %v", y0)
	}()

	go func() {
		y1 := x
		y1++

		log.Printf("y1: %v", y1)
	}()

	time.Sleep(1 * time.Second)
}
