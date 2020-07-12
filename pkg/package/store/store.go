package packagestore

import (
	"github.com/nikoksr/proji/pkg/domain"
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

func (ps *packageStore) LoadPackageList(labels ...string) ([]*domain.Package, error) {
	lenLabels := len(labels)
	if lenLabels < 1 {
		return ps.getAllPackages()
	}
	packages := make([]*domain.Package, 0, lenLabels)
	for _, label := range labels {
		pkg, err := ps.LoadPackage(label)
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
