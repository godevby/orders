package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

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

	log.Infow("startup", "main", "started", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Infow("shutdown", "main", "completed")

	// --------------------------------------------------------

	// if err := http.ListenAndServe(":7777", http.HandlerFunc(sayHello)); err != nil {
	// 	log.Fatalw("startup", "http server", err)
	// }

	api := http.Server{
		Addr:    ":7777",
		Handler: http.HandlerFunc(sayHello),
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// --------------------------------------------------------

	select {
	case err := <-serverErrors:
		log.Fatalw("startup", "http server", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Infow("shutdown", "graceful shutdown did not complete in time", err)
			err = api.Close()
		}
		if err != nil {
			log.Fatalw("shutdown", "could not stop http server gracefully", err)
		}
	}

	return nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Yey!", r.URL.Path)
}
