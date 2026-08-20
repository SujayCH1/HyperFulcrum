package connections

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"testing"
)

func newTestPool(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open(
		"postgres",
		"postgres://user:password@localhost:1/database?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func poolClosed(db *sql.DB) bool {
	err := db.Ping()
	return err != nil && strings.Contains(err.Error(), "database is closed")
}

func TestPoolStoreSetReplacesAndClosesPool(t *testing.T) {
	store := NewPoolStore()
	oldDB := newTestPool(t)
	newDB := newTestPool(t)
	defer store.Close()

	if err := store.Set("project-1", "node-1", oldDB); err != nil {
		t.Fatal(err)
	}

	if err := store.Set("project-1", "node-1", newDB); err != nil {
		t.Fatal(err)
	}

	if !poolClosed(oldDB) {
		t.Fatal("replaced pool was not closed")
	}

	db, err := store.Get("project-1", "node-1")
	if err != nil {
		t.Fatal(err)
	}
	if db != newDB {
		t.Fatal("new pool was not stored")
	}
}

func TestPoolStoreRejectsNilPool(t *testing.T) {
	store := NewPoolStore()

	if err := store.Set("project-1", "node-1", nil); !errors.Is(err, ErrNilPool) {
		t.Fatalf("expected nil pool error, got %v", err)
	}
}

func TestPoolStoreRemoveProject(t *testing.T) {
	store := NewPoolStore()
	db1 := newTestPool(t)
	db2 := newTestPool(t)

	if err := store.Set("project-1", "node-1", db1); err != nil {
		t.Fatal(err)
	}
	if err := store.Set("project-1", "node-2", db2); err != nil {
		t.Fatal(err)
	}

	if err := store.RemoveProject("project-1"); err != nil {
		t.Fatal(err)
	}

	if !poolClosed(db1) || !poolClosed(db2) {
		t.Fatal("project pools were not closed")
	}

	if _, err := store.Get("project-1", "node-1"); !errors.Is(err, ErrProjectPoolsNotFound) {
		t.Fatalf("expected project pools not found, got %v", err)
	}
}

func TestPoolStoreConcurrentAccess(t *testing.T) {
	store := NewPoolStore()
	defer store.Close()

	pools := make([]*sql.DB, 20)
	for index := range pools {
		pools[index] = newTestPool(t)
	}

	var waitGroup sync.WaitGroup

	for _, db := range pools {
		waitGroup.Add(1)

		go func(db *sql.DB) {
			defer waitGroup.Done()

			if err := store.Set("project-1", "node-1", db); err != nil {
				t.Error(err)
			}

			_, _ = store.Get("project-1", "node-1")
		}(db)
	}

	waitGroup.Wait()
}
