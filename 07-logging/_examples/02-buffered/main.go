package main

import (
	"bufio"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	w := bufio.NewWriter(os.Stderr)
	// defer w.Flush()
	logger := log.NewLogfmtLogger(w)

	for i := 1000; i > 0; i-- {
		logger.Log("msg", "countdown", "i", i)
	}

	// risky()
}

func risky() {
	os.Exit(1)
}
