// This program adds cancellation to a client request. Timeouts can be handled
// by setting the Timeout field on the client or by cancelling a Context
// associated with the request.

package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	c := http.Client{
		/*Timeout: 5 * time.Second,*/
	}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	// Make a new cancellable context with a timeout that derives from the
	// request's Context.
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(ctx) // Note assignment

	log.Println("sending request")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("executing request: %s", err)
	}
	defer resp.Body.Close()

	log.Printf("server responded with code %v", resp.StatusCode)
}
