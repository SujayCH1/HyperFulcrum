package hyperfulcrum

import (
	"context"
	App "hyperfulcrum/internal/app"
	"hyperfulcrum/pkg/logger"
)

func main() {
	ctx := context.Background()

	// initilize application instance
	application, err := App.NewApp(ctx)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application", "error", err)
	}

	// initialize application services
	application, err = application.Start(ctx, application)
	if err != nil {
		logger.Logger.Error("Failed to initizlize application services", "error", err)
	}
}
