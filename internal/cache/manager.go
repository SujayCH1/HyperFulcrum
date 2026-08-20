package cache

type CacheManager struct {
	Projects    *ProjectStore
	Nodes       *NodeStore
	Connections *ConnectionStore
	Topology    *NodeTopologyStore

	Columns       *ColumnStore
	FKEdges       *FKEdgesStore
	SchemaVersion *SchemaVersionStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectStore(),
		Nodes:       NewNodeStore(),
		Connections: NewConnectionStore(),
		Topology:    NewTopologyStore(),

		Columns:       NewColumnStore(),
		FKEdges:       NewFKEdgesStore(),
		SchemaVersion: NewSchemaVersionStore(),
	}
}

func (m *CacheManager) DeleteProject(projectID string) {
	m.Projects.Delete(projectID)
	m.Nodes.DeleteProject(projectID)
	m.Connections.DeleteProject(projectID)
	m.Topology.DeleteByProjectID(projectID)
	m.Columns.DeleteProject(projectID)
	m.FKEdges.DeleteProject(projectID)
	m.SchemaVersion.DeleteProject(projectID)
}

func (m *CacheManager) Clear() {
	m.Projects.Clear()
	m.Nodes.Clear()
	m.Connections.Clear()
	m.Topology.Clear()
	m.Columns.Clear()
	m.FKEdges.Clear()
	m.SchemaVersion.Clear()
}
