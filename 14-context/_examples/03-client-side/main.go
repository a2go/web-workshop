// This program shows an HTTP client making a request. By default there is no
// timeout or cancellation that happens from the client's end.

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

	log.Println("sending request")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("executing request: %s", err)
	}
	defer resp.Body.Close()

	log.Printf("server responded with code %v", resp.StatusCode)
}
