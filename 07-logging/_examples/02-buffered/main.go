// This program shows one kind of structured logger using go-kit.
// Note that the io.Writer in use is buffered so you must ensure it gets a
// chance to flush all of its output.
//
// - Run the program once and see the output is interrupted.
// - Uncomment the defer and run it again.
// - Uncomment the call to risky() and see it is interrupted again as os.Exit does not respect defers.

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
