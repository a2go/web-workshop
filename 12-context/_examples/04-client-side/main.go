package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	c := http.Client{
		/*Timeout: time.Minute,*/
	}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	// Add context to request.
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	req = req.WithContext(ctx) // Note assignment

	log.Print("sending request")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("executing request: %s", err)
	}
	defer resp.Body.Close()

	log.Printf("server responded with code %v", resp.StatusCode)
}
