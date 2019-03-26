// This program shows a server that uses an incoming request's context to
// control cancellation for a long operation like an outgoing request to
// another server.

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("start")
		defer log.Println("end")

		if err := getSomething(r.Context(), "http://localhost:8080"); err != nil {
			log.Printf("doRequest: %s", err)
		}
	})

	addr := ":8081"
	log.Printf("listening on address %q", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getSomething(ctx context.Context, addr string) error {
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return errors.Wrap(err, "making request")
	}
	req = req.WithContext(ctx)

	log.Printf("sending request to %q", addr)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	log.Printf("server responded with code %v", resp.StatusCode)

	return nil
}
