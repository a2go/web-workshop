// This program shows a simple server that launches a long-running operation
// with no cancellation. Run a "curl" to this server and quickly abort it. See
// what happens.

package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("start")
		defer log.Print("end")

		longOperation()
	})

	addr := ":8080"
	log.Printf("listening on address %q", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func longOperation() {
	for i, n := 1, 10; i <= n; i++ {
		time.Sleep(time.Second)

		log.Printf("%v/%v", i, n)
	}
}
