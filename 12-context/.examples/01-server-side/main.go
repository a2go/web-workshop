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

	log.Print(http.ListenAndServe(":8080", nil))
}

func longOperation() {
	for i, n := 1, 10; i <= n; i++ {
		time.Sleep(time.Second)

		log.Printf("%v percent complete", i*n)
	}
}
