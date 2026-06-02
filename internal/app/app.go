// internal/app/app.go
package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"hyperfulcrum/internal/api/router"
	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/connections"
	"hyperfulcrum/internal/database"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/pkg/logger"
)

type Application struct {
	Context context.Context

	// database pool
	database *sql.DB

	// repositories
	ProjectRepo  *repository.ProjectRepository
	NodeRepo     *repository.NodeRepository
	NodeConnRepo *repository.NodeConnectionRepository

	// connection pool
	ConnectionStore   *connections.ConnectionStore
	ConnectionManager *connections.ConnectionManager

	// cache
	CacheManager *cache.CacheManager
	CacheFetcher *cache.Fetcher

	// api
	Server *http.ServeMux
}

func New(ctx context.Context) (*Application, error) {
	return &Application{
		Context: ctx,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	logger.Logger.Info("Starting HyperFulcrum application")

	// Database
	logger.Logger.Info("Connecting to application database")
	db, err := database.NewDatabasePool(ctx)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	a.database = db
	logger.Logger.Info("Database connection established")

	// Migrations
	logger.Logger.Info("Running database migrations")
	if err := database.RunMigrations(); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}
	logger.Logger.Info("Database migrations completed")

	// Repositories
	logger.Logger.Info("Initializing repositories")

	a.ProjectRepo = repository.NewProjectRepository(a.database)
	a.NodeRepo = repository.NewNodeRepository(a.database)
	a.NodeConnRepo = repository.NewNodeConnectionRepository(a.database)

	logger.Logger.Info("Repositories initialized")

	// Cache
	logger.Logger.Info("Initializing cache manager")

	a.CacheManager = cache.NewCacheManager()
	a.CacheFetcher = cache.NewFetcher(
		*a.ProjectRepo,
		*a.NodeRepo,
		*a.NodeConnRepo,
		a.CacheManager,
	)

	logger.Logger.Info("Cache manager initialized")

	// Connections
	logger.Logger.Info("Initializing connection manager")

	a.ConnectionStore = connections.NewConnectionStore()
	a.ConnectionManager = connections.NewConnectionManager(
		a.ConnectionStore,
		a.ProjectRepo,
		a.NodeRepo,
		a.NodeConnRepo,
	)

	logger.Logger.Info("Connection manager initialized")

	//API Server
	logger.Logger.Info("Initializing application server")

	a.Server = router.NewRouter()

	logger.Logger.Info("Appllication server initialized")

	// application setup complete

	logger.Logger.Info("HyperFulcrum application started successfully")

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	logger.Logger.Info("Application shutdown initiated")

	if a.database != nil {
		logger.Logger.Info("Closing database connection")

		if err := a.database.Close(); err != nil {
			logger.Logger.Error("Failed to close database connection", "error", err)
			return err
		}

		logger.Logger.Info("Database connection closed")
	}

	logger.Logger.Info("HyperFulcrum application shutdown complete")

	return nil
}
