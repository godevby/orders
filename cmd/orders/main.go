package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	if err := run(logger); err != nil {
		logger.Error("startup", zap.Error(err))
		logger.Sync()
		os.Exit(1)
	}
}

func run(logger *zap.Logger) error {
	log := logger.Sugar()

	// --------------------------------------------------------

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// --------------------------------------------------------

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown completed", "signal", sig)

	return nil
}
