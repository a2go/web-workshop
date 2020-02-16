// This program shows adds cancellation to a long-running operation. Start a
// request and quickly cancel it. Notice the longOperation terminates before
// going all the way to 10.

package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("start")
		defer log.Println("end")

		longOperation(r.Context())
	})

	addr := ":8080"
	log.Printf("listening on address %q", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func longOperation(ctx context.Context) {
	for i, n := 1, 10; i <= n; i++ {
		select {
		case <-ctx.Done():
			log.Println("context was canceled: abort")
			return
		default:
		}

		time.Sleep(time.Second)

		log.Printf("%v/%v", i, n)
	}
}
