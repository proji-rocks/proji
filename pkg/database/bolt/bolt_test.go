package bolt

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TODO: Think about using filesystem mock instead of real filesystem.
func TestConnect(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "proji_test_*")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// Connect to the database
	dbPath := filepath.Join(dir, "test.db")
	db, err := Connect(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("bolt.Connect() failed: %v", err)
	}
	if db == nil {
		t.Fatalf("bolt.Connect() returned nil")
	}
	if db.Core == nil {
		t.Fatalf("bolt.Connect() returned a DB with a nil Core")
	}

	// Close the database
	if err = db.Close(context.Background()); err != nil {
		t.Fatalf("bolt.Close() failed: %v", err)
	}
}
