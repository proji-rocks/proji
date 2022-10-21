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
	// ErrProjectNotFound is returned when a project is not found in the repository.
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectExists is returned when a project with the same id already exists in the repository.
	ErrProjectExists = errors.New("project already exists")
)

type projectRepo struct {
	db         *bolt.DB
	bucketName string
}

// Compile-time check to ensure that projectRepo implements the domain.ProjectRepo interface.
var _ domain.ProjectRepo = (*projectRepo)(nil)

// New returns a new instance of the project repository. It requires a bolt database.
func New(db *db.DB) (domain.ProjectRepo, error) {
	if db == nil || db.Core == nil {
		return nil, errors.New("database is nil")
	}

	return &projectRepo{
		db:         db.Core,
		bucketName: domain.Project{}.Bucket(),
	}, nil
}

// Fetch fetches all projects from the database.
func (p projectRepo) Fetch(ctx context.Context) ([]domain.Project, error) {
	// Call projects from the database.
	var projects []domain.Project
	err := p.db.View(func(tx *bolt.Tx) error {
		// Open the bucket.
		bucket := tx.Bucket([]byte(p.bucketName))
		if bucket == nil {
			return db.ErrBucketNotFound
		}

		// Pre-alloc the project list.
		projects = make([]domain.Project, 0, bucket.Stats().KeyN)

		// Iterate over the bucket.
		return bucket.ForEach(func(_, projectData []byte) error {
			// Check if context is canceled.
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Unmarshal the project.
			project := domain.Project{}
			if err := json.Unmarshal(projectData, &project); err != nil {
				return errors.Wrap(err, "unmarshal project")
			}

			// Add the project to the list.
			projects = append(projects, project)

			return nil
		})
	})

	return projects, err
}

// GetByID fetches a project from the database by id.
func (p projectRepo) GetByID(ctx context.Context, id string) (domain.Project, error) {
	// Call project from database by its id.
	var project domain.Project
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

		// Get the project data.
		projectData := bucket.Get([]byte(id))
		if projectData == nil {
			return ErrProjectNotFound
		}

		// Unmarshal the project.
		if err := json.Unmarshal(projectData, &project); err != nil {
			return errors.Wrap(err, "unmarshal project")
		}

		return nil
	})

	return project, err
}

func (p projectRepo) doesProjectExist(ctx context.Context, path string) bool {
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

		// Get the project data.
		if bucket.Get([]byte(path)) == nil {
			return ErrProjectNotFound
		}

		return nil
	})

	return err == nil
}

// Store stores a project in the database.
func (p projectRepo) Store(ctx context.Context, project *domain.ProjectAdd) error {
	// Store the project in the database.
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

		// Marshal the project.
		projectData, err := json.Marshal(project)
		if err != nil {
			return errors.Wrap(err, "marshal project")
		}

		// Check if project already exists.
		if p.doesProjectExist(ctx, project.Path) {
			return ErrProjectExists
		}

		// Store the project.
		if err = bucket.Put([]byte(project.Path), projectData); err != nil {
			return errors.Wrap(err, "store project")
		}

		return nil
	})
}

// Update updates a project in the database.
func (p projectRepo) Update(ctx context.Context, project *domain.ProjectUpdate) error {
	// Update the project in the database.
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

		// Check if project already exists.
		if !p.doesProjectExist(ctx, project.ID) {
			return ErrProjectNotFound
		}

		// Marshal the project.
		projectData, err := json.Marshal(project)
		if err != nil {
			return errors.Wrap(err, "marshal project")
		}

		// Store/update the project. This will overwrite the existing project. Comparable to a PUT.
		if err = bucket.Put([]byte(project.ID), projectData); err != nil {
			return errors.Wrap(err, "store project")
		}

		return nil
	})
}

// Remove removes a project from the database.
func (p projectRepo) Remove(ctx context.Context, id string) error {
	// Remove the project from the database.
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

		// Check if project exists.
		if !p.doesProjectExist(ctx, id) {
			return ErrProjectNotFound
		}

		// Remove the project.
		if err := bucket.Delete([]byte(id)); err != nil {
			return errors.Wrap(err, "remove project")
		}

		return nil
	})
}
