package packageservice

import (
	"fmt"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/domain"
)

type packageService struct {
	authentication *config.APIAuthentication
	packageStore   domain.PackageStore
}

func New(auth *config.APIAuthentication, store domain.PackageStore) domain.PackageService {
	return &packageService{
		authentication: auth,
		packageStore:   store,
	}
}

func (ps packageService) StorePackage(pkg *domain.Package) error {
	if pkg == nil {
		return fmt.Errorf("received nil package")
	}

	return ps.packageStore.StorePackage(pkg)
}

func (ps packageService) LoadPackage(label string) (*domain.Package, error) {
	return ps.packageStore.LoadPackage(label)
}

func (ps packageService) LoadPackageList(labels ...string) ([]*domain.Package, error) {
	return ps.packageStore.LoadPackageList(labels...)
}

func (ps packageService) RemovePackage(label string) error {
	return ps.packageStore.RemovePackage(label)
}

func (ps packageService) PurgePackage(label string) error {
	return ps.packageStore.PurgePackage(label)
}
