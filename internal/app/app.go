// internal/app/app.go
package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hyperfulcrum/internal/api/handlers"
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/api/router"
	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/connections"
	"hyperfulcrum/internal/database"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/replication"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/internal/schema"
	"hyperfulcrum/internal/shardkey"
	"hyperfulcrum/pkg/logger"
	"net/http"
	"time"
)

type Application struct {
	Context context.Context

	// database pool
	database *sql.DB

	// repositories
	ProjectRepo       *repository.ProjectRepository
	NodeRepo          *repository.NodeRepository
	NodeConnRepo      *repository.NodeConnectionRepository
	TopologyRepo      *repository.NodeTopologyRepository
	ColumnRepo        *repository.ColumnRepository
	FKEdgesRepo       *repository.FKEdgesRepository
	ShardKeyRepo      *repository.ShardKeysRepository
	SchemaVersionRepo *repository.SchemaVersionRepository

	// metadata services
	ProjectService        *metadata.ProjectService
	NodeService           *metadata.NodeService
	NodeConnectionService *metadata.NodeConnectionService
	NodeTopologyService   *metadata.TopologyService
	ColumnService         *metadata.ColumnService
	FKEdgesService        *metadata.FKEdgesService
	SchemaVersionService  *metadata.SchemaVersionService

	// replication services
	ReplicationService *replication.ReplicationService
	InferenceService *shardkey.InferenceService

	// schema service
	SchemaService *schema.SchemaService

	// handlers
	ProjectHandler        *handlers.ProjectHandler
	NodeHandler           *handlers.NodeHandler
	NodeConnectionHandler *handlers.NodeConnectionHandler
	NodeTopologyHandler   *handlers.TopologyHandler

	// replication handlers
	ReplicationHandler *handlers.ReplicationHandler

	// connection pool
	PoolStore         *connections.PoolStore
	ConnectionManager *connections.ConnectionManager

	// cache
	CacheManager   *cache.CacheManager
	CacheRefresher *cache.CacheRefresher

	// api
	Server *http.Server
}

func New(ctx context.Context) (*Application, error) {
	return &Application{
		Context: ctx,
	}, nil
}

func (a *Application) Start(ctx context.Context) (startErr error) {
	logger.Logger.Info("Starting HyperFulcrum application")

	defer func() {
		if startErr == nil {
			return
		}

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		startErr = errors.Join(startErr, a.Stop(cleanupCtx))
	}()

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
	a.ShardKeyRepo = repository.NewShardKeysRepository(a.database)
	a.TopologyRepo = repository.NewNodeTopologyRepo(a.database)
	a.ColumnRepo = repository.NewColumnRepository(a.database)
	a.FKEdgesRepo = repository.NewFKEdgesRepository(a.database)
	a.SchemaVersionRepo = repository.NewSchemaVersionRepository(a.database)

	logger.Logger.Info("Repositories initialized")

	// Cache
	logger.Logger.Info("Initializing cache manager")

	a.CacheManager = cache.NewCacheManager()

	a.CacheRefresher = cache.NewCacheRefresher(
		a.ProjectRepo,
		a.NodeRepo,
		a.NodeConnRepo,
		a.TopologyRepo,
		a.ColumnRepo,
		a.FKEdgesRepo,
		a.SchemaVersionRepo,
		a.CacheManager,
	)

	logger.Logger.Info("Cache manager initialized")

	logger.Logger.Info("Loading metadata into cache")

	if err := a.CacheRefresher.RefreshAllProjects(ctx); err != nil {
		return fmt.Errorf("load metadata cache: %w", err)
	}

	logger.Logger.Info("Metadata cache loaded")

	// Connections
	logger.Logger.Info("Initializing connection manager")

	a.PoolStore = connections.NewPoolStore()
	a.ConnectionManager = connections.NewConnectionManager(
		a.PoolStore,
		a.ProjectRepo,
		a.NodeRepo,
		a.NodeConnRepo,
	)

	logger.Logger.Info("Connection manager initialized")

	if err := a.ConnectionManager.InitializeActiveConnections(ctx); err != nil {
		return fmt.Errorf("initialize active connections: %w", err)
	}

	logger.Logger.Info("Connected to all active shards")

	//API Server
	logger.Logger.Info("Initializing application server")

	// metadata services
	logger.Logger.Info("Initializing metadata services")

	a.ProjectService = metadata.NewProjectService(
		a.ProjectRepo,
		a.CacheManager,
		a.CacheRefresher,
		a.ConnectionManager,
	)

	a.NodeService = metadata.NewNodeService(
		a.NodeRepo,
		a.CacheManager,
		a.CacheRefresher,
		a.ConnectionManager,
	)

	a.NodeConnectionService = metadata.NewNodeConnectionService(
		a.NodeConnRepo,
		a.NodeRepo,
		a.CacheManager,
		a.CacheRefresher,
		a.ConnectionManager,
	)

	a.NodeTopologyService = metadata.NewTpologyService(
		a.TopologyRepo,
		a.CacheManager,
		a.CacheRefresher,
	)

	a.ColumnService = metadata.NewColumnService(
		a.ColumnRepo,
		a.CacheManager,
		a.CacheRefresher,
	)

	a.FKEdgesService = metadata.NewFKEdgesService(
		a.FKEdgesRepo,
		a.CacheManager,
		a.CacheRefresher,
	)

	a.SchemaVersionService = metadata.NewSchemaVersionService(
		a.SchemaVersionRepo,
		a.CacheManager,
		a.CacheRefresher,
	)

	logger.Logger.Info("Initialized metadata services")

	// MAJOR SERVICES
	// replication services
	logger.Logger.Info("Initializing replication service")

	a.ReplicationService = replication.NewReplicationService(
		a.NodeTopologyService,
		a.NodeService,
	)

	logger.Logger.Info("Replication service initialized")

	//schema service
	logger.Logger.Info("Initializing schema service")

	a.SchemaService = schema.NewSchemaService(
		a.ColumnService,
		a.FKEdgesService,
		a.SchemaVersionService,
	)

	logger.Logger.Info("Schema service initialized")
   
	logger.Logger.Info("Initializing shardkey inference service")
	a.InferenceService = shardkey.NewInferenceService(
		a.ColumnRepo,
		a.FKEdgesRepo,
		a.ShardKeyRepo,
	)
	logger.Logger.Info("Shardkey inference service initialized")

	// handlers
	a.ProjectHandler = handlers.NewProjectHandler(
		a.ProjectService,
	)

	a.NodeHandler = handlers.NewNodeHandler(
		a.NodeService,
	)

	a.NodeConnectionHandler = handlers.NewNodeConnectionHandler(
		a.NodeConnectionService,
	)

	a.NodeTopologyHandler = handlers.NewTopoogyHandler(
		a.NodeTopologyService,
	)

	a.ReplicationHandler = handlers.NewReplicationHandler(
		a.ReplicationService,
	)

	// server setup
	mux := router.NewRouter(
		a.ProjectHandler,
		a.NodeHandler,
		a.NodeConnectionHandler,
		a.NodeTopologyHandler,
		a.ReplicationHandler,
	)

	handler := middleware.RequestID(
		middleware.RequestLogger(
			middleware.Recovery(
				middleware.RequestTimeout(
					30*time.Second,
					middleware.CORS(mux),
				),
			),
		),
	)

	a.Server = &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      35 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	logger.Logger.Info("Application server initialized")

	go func() {
		logger.Logger.Info("Starting HTTP server on :8080")
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("HTTP server failed", "error", err)
		}
	}()

	// application setup complete

	logger.Logger.Info("HyperFulcrum application started successfully")

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	logger.Logger.Info("Application shutdown initiated")

	var stopErr error

	if a.Server != nil {
		logger.Logger.Info("Stopping HTTP server")

		if err := a.Server.Shutdown(ctx); err != nil {
			logger.Logger.Error("Failed to stop HTTP server", "error", err)
			stopErr = errors.Join(stopErr, err)
		} else {
			logger.Logger.Info("HTTP server stopped")
		}
	}

	if a.ConnectionManager != nil {
		logger.Logger.Info("Closing node connection pools")

		if err := a.ConnectionManager.Close(); err != nil {
			logger.Logger.Error("Failed to close node connection pools", "error", err)
			stopErr = errors.Join(stopErr, err)
		} else {
			logger.Logger.Info("Node connection pools closed")
		}
	}

	if a.database != nil {
		logger.Logger.Info("Closing database connection")

		if err := a.database.Close(); err != nil {
			logger.Logger.Error("Failed to close database connection", "error", err)
			stopErr = errors.Join(stopErr, err)
		} else {
			logger.Logger.Info("Database connection closed")
		}
	}

	logger.Logger.Info("HyperFulcrum application shutdown complete")

	return stopErr
}
