package main

import (
	"log"
	"net/http"
)

func main() {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	log.Print("sending request")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("executing request: %s", err)
	}
	defer resp.Body.Close()

	log.Printf("server responded with code %v", resp.StatusCode)
}
