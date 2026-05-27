package main

import (
	"context"
	App "hyperfulcrum/internal/app"
	"hyperfulcrum/pkg/logger"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initilize application instance
	application, err := App.New(ctx)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application", "error", err)
	}

	// initialize application services
	err = application.Start(ctx)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application services", "error", err)
	}

	<-ctx.Done()

	// stop the application
	err = application.Stop(context.Background())
	if err != nil {
		logger.Logger.Error("Failed to stop application", "error", err)
	}

}
