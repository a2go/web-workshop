package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("server listening on", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-osSignals:
		log.Println("caught signal, shutting down")

		// Give outstanding requests 5 seconds to complete.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: gracefully shutting down server: %s", err)
			if err := server.Close(); err != nil {
				log.Printf("error: closing server: %s", err)
			}
		}
	}

	log.Println("done")
}

// Echo is a basic HTTP Handler.
func Echo(w http.ResponseWriter, r *http.Request) {

	// Print a random number at the beginning and end of each request.
	n := rand.Intn(1000)
	log.Println("start", n)
	defer log.Println("end", n)

	// Simulate a long-running request.
	time.Sleep(3 * time.Second)

	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
}
