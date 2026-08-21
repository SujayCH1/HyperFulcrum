package main

import (
	"context"
	App "hyperfulcrum/internal/app"
	"hyperfulcrum/pkg/logger"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initilize application instance
	application, err := App.New(ctx)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application", "error", err)
		return
	}

	// initialize application services
	err = application.Start(ctx)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application services", "error", err)
		return
	}

	<-ctx.Done()

	// stop the application
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = application.Stop(shutdownCtx)
	if err != nil {
		logger.Logger.Error("Failed to stop application", "error", err)
	}

}
