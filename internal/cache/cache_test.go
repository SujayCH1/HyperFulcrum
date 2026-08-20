package cache

import (
	"sync"
	"testing"

	"hyperfulcrum/internal/repository"
)

func TestProjectStoreTracksEmptyLoadedState(t *testing.T) {
	store := NewProjectStore()

	projects, loaded := store.GetAll()
	if loaded || len(projects) != 0 {
		t.Fatal("new project store must be unloaded and empty")
	}

	store.Replace([]repository.Project{})

	projects, loaded = store.GetAll()
	if !loaded || len(projects) != 0 {
		t.Fatal("empty project replacement must be cached as loaded")
	}
}

func TestNodeStoreReplacesProjectAtomically(t *testing.T) {
	store := NewNodeStore()

	first := []repository.Node{
		{ID: "node-1", ProjectID: "wrong-project"},
	}
	second := []repository.Node{
		{ID: "node-1", ProjectID: "project-1"},
		{ID: "node-2", ProjectID: "project-1"},
	}

	store.ReplaceProject("project-1", first)

	var waitGroup sync.WaitGroup
	errors := make(chan string, 4)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for index := 0; index < 1000; index++ {
			if index%2 == 0 {
				store.ReplaceProject("project-1", first)
			} else {
				store.ReplaceProject("project-1", second)
			}
		}
	}()

	for reader := 0; reader < 4; reader++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			for index := 0; index < 1000; index++ {
				nodes, loaded := store.GetByProject("project-1")
				if !loaded || len(nodes) < 1 || len(nodes) > 2 {
					errors <- "project replacement exposed a partial state"
					return
				}
			}
		}()
	}

	waitGroup.Wait()
	close(errors)

	for err := range errors {
		t.Fatal(err)
	}

	store.ReplaceProject("project-1", nil)

	nodes, loaded := store.GetByProject("project-1")
	if !loaded || len(nodes) != 0 {
		t.Fatal("empty node replacement must be cached as loaded")
	}

	if _, ok := store.Get("node-1"); ok {
		t.Fatal("node index must remove stale entries")
	}
}

func TestColumnStoreUsesTypedKeys(t *testing.T) {
	store := NewColumnStore()

	columns := []repository.Column{
		{ProjectID: "project-1", TableName: "a:b", ColumnName: "c"},
		{ProjectID: "project-1", TableName: "a", ColumnName: "b:c"},
	}

	store.ReplaceProject("project-1", columns)

	if _, ok := store.Get("project-1", "a:b", "c"); !ok {
		t.Fatal("first column key was overwritten")
	}

	if _, ok := store.Get("project-1", "a", "b:c"); !ok {
		t.Fatal("second column key was overwritten")
	}
}

func TestFKEdgesStoreUsesTypedKeys(t *testing.T) {
	store := NewFKEdgesStore()

	edges := []repository.FkEdges{
		{
			ProjectId:    "project-1",
			ParentTable:  "a:b",
			ParentColumn: "c",
			ChildTable:   "d",
			ChildColumn:  "e",
		},
		{
			ProjectId:    "project-1",
			ParentTable:  "a",
			ParentColumn: "b:c",
			ChildTable:   "d",
			ChildColumn:  "e",
		},
	}

	store.ReplaceProject("project-1", edges)

	if _, ok := store.Get("project-1", "a:b", "c", "d", "e"); !ok {
		t.Fatal("first foreign key edge was overwritten")
	}

	if _, ok := store.Get("project-1", "a", "b:c", "d", "e"); !ok {
		t.Fatal("second foreign key edge was overwritten")
	}
}

func TestSchemaVersionStoreTracksMissingValues(t *testing.T) {
	store := NewSchemaVersionStore()

	store.SetMissing("project-1")

	if !store.Loaded("project-1") {
		t.Fatal("missing schema version must be cached as loaded")
	}

	if _, ok := store.Get("project-1"); ok {
		t.Fatal("missing schema version must not create a value")
	}

	store.Set(repository.SchemaVersion{
		ID:        "schema-1",
		ProjectID: "project-1",
	})

	if _, ok := store.Get("project-1"); !ok {
		t.Fatal("schema version must replace the missing state")
	}
}

func TestCacheManagerDeletesProjectMetadata(t *testing.T) {
	manager := NewCacheManager()

	manager.Projects.Replace([]repository.Project{{ID: "project-1"}})
	manager.Nodes.ReplaceProject("project-1", []repository.Node{{ID: "node-1"}})
	manager.Connections.ReplaceProject(
		"project-1",
		[]repository.NodeConnection{{NodeId: "node-1"}},
	)
	manager.Topology.Set("project-1", []repository.NodeTopology{})
	manager.Columns.ReplaceProject("project-1", []repository.Column{})
	manager.FKEdges.ReplaceProject("project-1", []repository.FkEdges{})
	manager.SchemaVersion.SetMissing("project-1")

	manager.DeleteProject("project-1")

	if manager.Projects.Exists("project-1") {
		t.Fatal("project cache entry was not deleted")
	}

	if _, loaded := manager.Nodes.GetByProject("project-1"); loaded {
		t.Fatal("node cache state was not deleted")
	}

	if manager.Connections.ProjectLoaded("project-1") {
		t.Fatal("connection cache state was not deleted")
	}

	if manager.Topology.ExistsByProjectID("project-1") {
		t.Fatal("topology cache state was not deleted")
	}

	if _, loaded := manager.Columns.GetByProject("project-1"); loaded {
		t.Fatal("column cache state was not deleted")
	}

	if _, loaded := manager.FKEdges.GetByProject("project-1"); loaded {
		t.Fatal("foreign key cache state was not deleted")
	}

	if manager.SchemaVersion.Loaded("project-1") {
		t.Fatal("schema version cache state was not deleted")
	}
}
