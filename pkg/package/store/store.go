package packagestore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type packageStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.PackageStore {
	return &packageStore{
		db: db,
	}
}

func (ps *packageStore) StorePackage(pkg *domain.Package) error {
	err := ps.db.First(pkg, "label = ?", pkg.Label).Error
	if err == nil {
		return &PackageExistsError{Label: pkg.Label}
	}
	if err == gorm.ErrRecordNotFound {
		return ps.db.Create(pkg).Error
	}
	return err
}

func (ps *packageStore) LoadPackage(label string) (*domain.Package, error) {
	var pkg domain.Package
	err := ps.db.Preload(clause.Associations).First(&pkg, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &PackageNotFoundError{Label: label}
	}
	return &pkg, err
}

func (ps packageStore) LoadPackage(loadDependencies bool, label string) (*domain.Package, error) {
	conditions := fmt.Sprintf("label = '%s'", label)
	if loadDependencies {
		conditions = fmt.Sprintf("packages.label = '%s'", label)
	}
	return ps.loadPackage(loadDependencies, conditions)
}

func (ps packageStore) loadPackage(loadDependencies bool, conditions string) (*domain.Package, error) {
	if loadDependencies {
		return ps.deepQueryPackage(conditions)
	}
	return ps.queryPackage(conditions)
}

func (ps packageStore) LoadPackageList(loadDependencies bool, labels ...string) ([]*domain.Package, error) {
	var err error
	labelCount := len(labels)
	if labelCount < 1 {
		labels, err = ps.queryAllLabels()
		labelCount = len(labels)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// loadAllPackages loads and returns all packages found in the database.
func (ps *packageStore) getAllPackages() ([]*domain.Package, error) {
	var packages []*domain.Package
	err := ps.db.Preload(clause.Associations).Find(&packages).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &NoPackagesFoundError{}
	}
	return packages, err
}

func (ps *packageStore) RemovePackage(label string) error {
	err := ps.db.Delete(&domain.Package{}, "label = ? AND deleted_at IS NULL", label).Error
	if err == gorm.ErrRecordNotFound {
		return &PackageNotFoundError{Label: label}
	}
	return err
}

func (ps *packageStore) PurgePackage(label string) error {
	err := ps.db.Unscoped().Delete(&domain.Package{}, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return &PackageNotFoundError{Label: label}
	}
	return err
}
