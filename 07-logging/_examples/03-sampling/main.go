// Some loggers are "sampling" loggers. These are optimized for speed over
// precision. Using one of these loggers comes with the acknowledgment that not
// every message will be logged.
//
// Run in development mode vs production mode and notice the difference.

package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	//logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logger.Sync()

	for i := 1000; i > 0; i-- {
		logger.Info("countdown", zap.Int("i", i))
	}
}
