// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"hyperfulcrum/internal/database"
	"hyperfulcrum/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	// context
	ctx context.Context

	// databasee pool
	db *pgxpool.Pool

	// api server
	// router
	// repositories
	// services
}

func New(ctx context.Context) (*Application, error) {
	return &Application{
		ctx: ctx,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {

	// Connect to the databse
	logger.Logger.Info("Connecting to application database")
	pool, err := database.NewDatabasePool(ctx)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	a.db = pool

	// Run database migrations
	logger.Logger.Info("Running databse migrations")
	err = database.RunMigrations()
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	// shut down db
	logger.Logger.Info("Shutting down database")
	if a.db != nil {
		a.db.Close()
	}

	return nil
}
