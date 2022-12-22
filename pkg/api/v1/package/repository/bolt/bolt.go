package bolt

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	bolt "go.etcd.io/bbolt"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	db "github.com/nikoksr/proji/pkg/database/bolt"
)

var (
	// ErrPackageNotFound is returned when a package is not found in the repository.
	ErrPackageNotFound = errors.New("package not found")
	// ErrPackageExists is returned when a package with the same label already exists in the repository.
	ErrPackageExists = errors.New("package already exists")
)

type packageRepo struct {
	db         *bolt.DB
	bucketName string
}

// Compile-time check to ensure that packageRepo{} implements the domain.PackageRepo interface.
var _ domain.PackageRepo = (*packageRepo)(nil)

// New returns a new instance of the package repository. It requires a bolt database.
func New(db *db.DB) (domain.PackageRepo, error) {
	if db == nil || db.Core == nil {
		return nil, errors.New("database is nil")
	}

	return &packageRepo{
		db:         db.Core,
		bucketName: (&domain.Package{}).Bucket(),
	}, nil
}

// Fetch fetches all packages from the database.
func (p packageRepo) Fetch(ctx context.Context) ([]domain.Package, error) {
	// Call packages from the database.
	var packages []domain.Package
	err := p.db.View(func(tx *bolt.Tx) error {
		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Pre-alloc the package list.
		packages = make([]domain.Package, 0, bucket.Stats().KeyN)

		// Iterate over the bucket.
		return bucket.ForEach(func(_, pkgData []byte) error {
			// Check if context is canceled.
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Unmarshal the package.
			_package := domain.Package{}
			if err := json.Unmarshal(pkgData, &_package); err != nil {
				return errors.Wrap(err, "unmarshal package")
			}

			// Add the package to the list.
			packages = append(packages, _package)

			return nil
		})
	})

	return packages, err
}

// GetByLabel fetches a package from the database by label.
func (p packageRepo) GetByLabel(ctx context.Context, label string) (domain.Package, error) {
	// Call package from database by its label.
	var _package domain.Package
	err := p.db.View(func(tx *bolt.Tx) error {
		// Check if context is canceled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Get the package data.
		pkgData := bucket.Get([]byte(label))
		if pkgData == nil {
			return ErrPackageNotFound
		}

		// Unmarshal the package.
		if err := json.Unmarshal(pkgData, &_package); err != nil {
			return errors.Wrap(err, "unmarshal package")
		}

		return nil
	})

	return _package, err
}

func (p packageRepo) doesPackageExist(ctx context.Context, label string) bool {
	err := p.db.View(func(tx *bolt.Tx) error {
		// Check if context is canceled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Get the package data.
		if bucket.Get([]byte(label)) == nil {
			return ErrPackageNotFound
		}

		return nil
	})

	return err == nil
}

// Store stores a package in the database.
func (p packageRepo) Store(ctx context.Context, _package *domain.PackageAdd) error {
	// Store the package in the database.
	return p.db.Update(func(tx *bolt.Tx) error {
		// Check if context is canceled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Open the bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(p.bucketName))
		if err != nil {
			return errors.Wrap(err, "create bucket")
		}

		// Check if package with the same label already exists.
		if p.doesPackageExist(ctx, _package.Label) {
			return ErrPackageExists
		}

		// Marshal the package.
		pkgData, err := json.Marshal(_package)
		if err != nil {
			return errors.Wrap(err, "marshal package")
		}

		// Store the package.
		if err = bucket.Put([]byte(_package.Label), pkgData); err != nil {
			return errors.Wrap(err, "store package")
		}

		return nil
	})
}

// Update updates a package in the database.
func (p packageRepo) Update(ctx context.Context, _package *domain.PackageUpdate) error {
	// Update the package in the database.
	return p.db.Update(func(tx *bolt.Tx) error {
		// Check if context is canceled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Check if package exists.
		if !p.doesPackageExist(ctx, _package.Label) {
			return ErrPackageNotFound
		}

		// Marshal the package.
		pkgData, err := json.Marshal(_package)
		if err != nil {
			return errors.Wrap(err, "marshal package")
		}

		// Store/update the package. This will overwrite the existing package. Comparable to a PUT.
		if err = bucket.Put([]byte(_package.Label), pkgData); err != nil {
			return errors.Wrap(err, "store package")
		}

		return nil
	})
}

// Remove removes a package from the database.
func (p packageRepo) Remove(ctx context.Context, label string) error {
	// Remove the package from the database.
	return p.db.Update(func(tx *bolt.Tx) error {
		// Check if context is canceled.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Check if package exists.
		if !p.doesPackageExist(ctx, label) {
			return ErrPackageNotFound
		}

		// Remove the package.
		if err := bucket.Delete([]byte(label)); err != nil {
			return errors.Wrap(err, "remove package")
		}

		return nil
	})
}
