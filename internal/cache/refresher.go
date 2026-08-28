package cache

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/repository"
)

type CacheRefresher struct {
	projectRepo       *repository.ProjectRepository
	nodeRepo          *repository.NodeRepository
	connectionRepo    *repository.NodeConnectionRepository
	topologyRepo      *repository.NodeTopologyRepository
	shardRepo         *repository.ShardRepository
	runtimeRepo       *repository.NodeRuntimeStateRepository
	columnRepo        *repository.ColumnRepository
	fkRepo            *repository.FKEdgesRepository
	schemaVersionRepo *repository.SchemaVersionRepository
	shardKeyRepo      *repository.ShardKeyRepository

	cache *CacheManager
}

type projectMetadata struct {
	project       repository.Project
	nodes         []repository.Node
	connections   []repository.NodeConnection
	topologies    []repository.NodeTopology
	shards        []repository.Shard
	runtimeStates []repository.NodeRuntimeState
	columns       []repository.Column
	edges         []repository.FkEdges
	schemaVersion *repository.SchemaVersion
	shardKeys     []repository.ShardKey
}

func NewCacheRefresher(
	projectRepo *repository.ProjectRepository,
	nodeRepo *repository.NodeRepository,
	connectionRepo *repository.NodeConnectionRepository,
	topologyRepo *repository.NodeTopologyRepository,
	shardRepo *repository.ShardRepository,
	runtimeRepo *repository.NodeRuntimeStateRepository,
	columnRepo *repository.ColumnRepository,
	fkRepo *repository.FKEdgesRepository,
	schemaVersionRepo *repository.SchemaVersionRepository,
	shardKeyRepo *repository.ShardKeyRepository,
	cache *CacheManager,
) *CacheRefresher {

	return &CacheRefresher{
		projectRepo:       projectRepo,
		nodeRepo:          nodeRepo,
		connectionRepo:    connectionRepo,
		topologyRepo:      topologyRepo,
		shardRepo:         shardRepo,
		runtimeRepo:       runtimeRepo,
		columnRepo:        columnRepo,
		fkRepo:            fkRepo,
		schemaVersionRepo: schemaVersionRepo,
		shardKeyRepo:      shardKeyRepo,
		cache:             cache,
	}
}

func (r *CacheRefresher) RefreshProjectMetadata(
	ctx context.Context,
	projectID string,
) error {

	project, err := r.projectRepo.ProjectGetByID(
		ctx,
		projectID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.cache.DeleteProject(projectID)
			return nil
		}

		return err
	}

	metadata, err := r.fetchProjectMetadata(ctx, project)
	if err != nil {
		return err
	}

	r.storeProjectMetadata(metadata)

	return nil
}

func (r *CacheRefresher) RefreshProjects(
	ctx context.Context,
) error {

	projects, err := r.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	cachedProjects, loaded := r.cache.Projects.GetAll()
	if loaded {
		projectIDs := make(map[string]bool, len(projects))

		for _, project := range projects {
			projectIDs[project.ID] = true
		}

		for _, project := range cachedProjects {
			if !projectIDs[project.ID] {
				r.cache.DeleteProject(project.ID)
			}
		}
	}

	r.cache.Projects.Replace(projects)

	return nil
}

func (r *CacheRefresher) RefreshProject(
	ctx context.Context,
	projectID string,
) error {

	project, err := r.projectRepo.ProjectGetByID(
		ctx,
		projectID,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			r.cache.DeleteProject(projectID)
			return nil
		}

		return err
	}

	r.cache.Projects.Set(project)

	return nil
}

func (r *CacheRefresher) RefreshNodes(
	ctx context.Context,
	projectID string,
) error {

	nodes, err := r.nodeRepo.NodeList(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Nodes.ReplaceProject(projectID, nodes)

	return nil
}

func (r *CacheRefresher) RefreshConnections(
	ctx context.Context,
	projectID string,
) error {

	connections, err := r.connectionRepo.ConnectionsListByProjectID(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Connections.ReplaceProject(projectID, connections)

	return nil
}

func (r *CacheRefresher) RefreshTopology(
	ctx context.Context,
	projectID string,
) error {

	topologies, err := r.topologyRepo.TopologyGetAll(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Topology.Set(
		projectID,
		topologies,
	)

	return nil
}

func (r *CacheRefresher) RefreshShards(ctx context.Context, projectID string) error {
	shards, err := r.shardRepo.ShardList(ctx, projectID)
	if err != nil {
		return err
	}
	r.cache.Shards.ReplaceProject(projectID, shards)
	return nil
}

func (r *CacheRefresher) RefreshNodeRuntimeStates(ctx context.Context, projectID string) error {
	states, err := r.runtimeRepo.RuntimeStateListByProject(ctx, projectID)
	if err != nil {
		return err
	}
	r.cache.Runtime.ReplaceProject(projectID, states)
	return nil
}

func (r *CacheRefresher) RefreshAllProjects(
	ctx context.Context,
) error {

	projects, err := r.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	metadata := make([]projectMetadata, 0, len(projects))

	for _, project := range projects {
		snapshot, err := r.fetchProjectMetadata(ctx, project)
		if err != nil {
			return err
		}

		metadata = append(metadata, snapshot)
	}

	r.cache.Clear()
	r.cache.Projects.Replace(projects)

	for _, snapshot := range metadata {
		r.storeProjectMetadata(snapshot)
	}

	return nil
}

func (r *CacheRefresher) RefreshColumns(
	ctx context.Context,
	projectID string,
) error {

	columns, err := r.columnRepo.ColumnsListByProjectID(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Columns.ReplaceProject(projectID, columns)

	return nil
}

func (r *CacheRefresher) RefreshFKEdges(
	ctx context.Context,
	projectID string,
) error {

	edges, err := r.fkRepo.EdgesListByProjectID(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.FKEdges.ReplaceProject(projectID, edges)

	return nil
}

func (r *CacheRefresher) RefreshSchemaVersion(
	ctx context.Context,
	projectID string,
) error {

	schema, err := r.schemaVersionRepo.FetchSchema(
		ctx,
		projectID,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			r.cache.SchemaVersion.SetMissing(projectID)
			return nil
		}

		return err
	}

	r.cache.SchemaVersion.Set(schema)

	return nil
}

func (r *CacheRefresher) RefreshShardKeys(
	ctx context.Context,
	projectID string,
) error {

	keys, err := r.shardKeyRepo.FetchShardKeys(ctx, projectID)
	if err != nil {
		return err
	}

	r.cache.ShardKeys.ReplaceProject(projectID, keys)

	return nil
}

func (r *CacheRefresher) fetchProjectMetadata(
	ctx context.Context,
	project repository.Project,
) (projectMetadata, error) {

	nodes, err := r.nodeRepo.NodeList(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	connections, err := r.connectionRepo.ConnectionsListByProjectID(
		ctx,
		project.ID,
	)
	if err != nil {
		return projectMetadata{}, err
	}

	topologies, err := r.topologyRepo.TopologyGetAll(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	shards, err := r.shardRepo.ShardList(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	runtimeStates, err := r.runtimeRepo.RuntimeStateListByProject(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	columns, err := r.columnRepo.ColumnsListByProjectID(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	edges, err := r.fkRepo.EdgesListByProjectID(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	shardKeys, err := r.shardKeyRepo.FetchShardKeys(ctx, project.ID)
	if err != nil {
		return projectMetadata{}, err
	}

	schemaVersion, err := r.schemaVersionRepo.FetchSchema(ctx, project.ID)
	if err != nil && err != sql.ErrNoRows {
		return projectMetadata{}, err
	}

	metadata := projectMetadata{
		project:       project,
		nodes:         nodes,
		connections:   connections,
		topologies:    topologies,
		shards:        shards,
		runtimeStates: runtimeStates,
		columns:       columns,
		edges:         edges,
		shardKeys:     shardKeys,
	}

	if err == nil {
		metadata.schemaVersion = &schemaVersion
	}

	return metadata, nil
}

func (r *CacheRefresher) storeProjectMetadata(
	metadata projectMetadata,
) {

	projectID := metadata.project.ID

	r.cache.Projects.Set(metadata.project)
	r.cache.Nodes.ReplaceProject(projectID, metadata.nodes)
	r.cache.Connections.ReplaceProject(projectID, metadata.connections)
	r.cache.Topology.Set(projectID, metadata.topologies)
	r.cache.Shards.ReplaceProject(projectID, metadata.shards)
	r.cache.Runtime.ReplaceProject(projectID, metadata.runtimeStates)
	r.cache.Columns.ReplaceProject(projectID, metadata.columns)
	r.cache.FKEdges.ReplaceProject(projectID, metadata.edges)
	r.cache.ShardKeys.ReplaceProject(projectID, metadata.shardKeys)

	if metadata.schemaVersion == nil {
		r.cache.SchemaVersion.SetMissing(projectID)
		return
	}

	r.cache.SchemaVersion.Set(*metadata.schemaVersion)
}
