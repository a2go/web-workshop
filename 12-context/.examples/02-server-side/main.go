package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("start")
		defer log.Print("end")

		longOperation(r.Context())
	})

	log.Print(http.ListenAndServe(":8080", nil))
}

func longOperation(ctx context.Context) {
	for i, n := 1, 10; i <= n; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second)

		log.Printf("%v percent complete", i*n)
	}
}
