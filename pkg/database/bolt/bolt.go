package bolt

import (
	"context"
	"time"

	"github.com/nikoksr/simplog"
	bolt "go.etcd.io/bbolt"
)

// ErrBucketNotFound is returned when a bucket is not found. It's an alias for bolt.ErrBucketNotFound.
var ErrBucketNotFound = bolt.ErrBucketNotFound

// DB is a wrapper around a BoltDB connection. It is used to manage the database connection.
type DB struct {
	Core *bolt.DB
}

func connect(ctx context.Context, path string) (*DB, error) {
	logger := simplog.FromContext(ctx)

	// Check if file is already open.
	logger.Debugf("trying to open database file: %q", path)
	db, err := bolt.Open(path, 0o600, &bolt.Options{
		Timeout:    1 * time.Second,
		NoGrowSync: false,
		// FreelistType: bolt.FreelistArrayType,
		FreelistType: bolt.FreelistMapType,
	})

	return &DB{db}, err
}

// Connect connects to the database. It returns an error if the connection fails.
func Connect(ctx context.Context, path string) (*DB, error) {
	return connect(ctx, path)
}

// Close closes the database connection.
func (db *DB) Close(ctx context.Context) error {
	logger := simplog.FromContext(ctx)
	logger.Debugf("closing database connection")

	return db.Core.Close()
}
