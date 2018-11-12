package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	// logger, _ := zap.NewProduction()
	defer logger.Sync()

	for i := 1000; i > 0; i-- {
		logger.Info("countdown", zap.Int("i", i))
	}
}
