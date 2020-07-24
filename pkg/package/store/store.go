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

func (ps packageStore) StorePackage(pkg *domain.Package) error {
	// Check if package exists
	err := ps.db.Where("label = ?", pkg.Label).First(pkg).Error
	if err == nil {
		return ErrPackageExists
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	err = tx.Omit(clause.Associations).Create(pkg).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert package")
	}

	err = storeTemplates(tx, pkg.Templates, pkg.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = storePlugins(tx, pkg.Plugins, pkg.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func storeTemplates(tx *gorm.DB, templates []*domain.Template, packageID uint) error {
	var err error
	insertTemplateStmt := "INSERT OR IGNORE INTO templates (created_at, updated_at, is_file, destination, path, description) VALUES (?, ?, ?, ?, ?, ?)"
	insertAssociationStmt := "INSERT OR IGNORE INTO package_templates (package_id, template_id) VALUES (?, ?)"
	queryIDStmt := "SELECT id from templates WHERE destination = ? AND path = ?"
	for _, template := range templates {
		now := time.Now()
		err = tx.Exec(insertTemplateStmt, now, now, template.IsFile, template.Destination, template.Path, template.Description).Error
		if err != nil {
			return err
		}

		// TODO: Possibly improve. Maybe select can be done outside of loop once.
		rows, err := tx.Raw(queryIDStmt, template.Destination, template.Path).Rows()
		if err != nil {
			return err
		}
		if rows.Err() != nil {
			return err
		}
		for rows.Next() {
			var id null.Int
			err = rows.Scan(&id)
			if err != nil {
				return err
			}
			if !id.Valid {
				return err
			}
			template.ID = uint(id.Int64)
		}
		err = rows.Close()
		if err != nil {
			return err
		}

		err = tx.Exec(insertAssociationStmt, packageID, template.ID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func storePlugins(tx *gorm.DB, plugins []*domain.Plugin, packageID uint) error {
	var err error
	insertPluginStmt := "INSERT OR IGNORE INTO plugins (created_at, updated_at, path, exec_number, description) VALUES (?, ?, ?, ?, ?)"
	insertAssociationStmt := "INSERT OR IGNORE INTO package_plugins (package_id, plugin_id) VALUES (?, ?)"
	queryIDStmt := "SELECT id from plugins WHERE path = ?"
	for _, plugin := range plugins {
		now := time.Now()
		err = tx.Exec(insertPluginStmt, now, now, plugin.Path, plugin.ExecNumber, plugin.Description).Error
		if err != nil {
			return err
		}

		// TODO: Possibly improve. Maybe select can be done outside of loop once.
		rows, err := tx.Raw(queryIDStmt, plugin.Path).Rows()
		if err != nil {
			return err
		}
		if rows.Err() != nil {
			return err
		}
		for rows.Next() {
			var id null.Int
			err = rows.Scan(&id)
			if err != nil {
				return err
			}
			if !id.Valid {
				return err
			}
			plugin.ID = uint(id.Int64)
		}
		err = rows.Close()
		if err != nil {
			return err
		}

		err = tx.Exec(insertAssociationStmt, packageID, plugin.ID).Error
		if err != nil {
			return err
		}
	}
	return nil
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
	}
	packageList := make([]*domain.Package, 0, labelCount)
	mut := sync.Mutex{}
	errs, _ := errgroup.WithContext(context.Background())
	for _, label := range labels {
		label := label
		errs.Go(func() error {
			pkg, err := ps.LoadPackage(loadDependencies, label)
			if err != nil {
				return err
			}
			mut.Lock()
			packageList = append(packageList, pkg)
			mut.Unlock()
			return nil
		})
	}
	return packageList, errs.Wait()
}

const (
	defaultPackageQueryBase     = `SELECT name, label, description FROM packages`
	defaultPackageDeepQueryBase = `SELECT
	packages.name,
	packages.label,
	packages.description as package_description,
	templates.is_file,
	templates.destination,
	templates."path" as template_path,
	templates.description as template_description,
	plugins."path" as plugin_path,
	plugins.exec_number,
	plugins.description as plugin_description
	FROM packages
LEFT OUTER JOIN package_templates
	ON packages.id = package_templates.package_id
LEFT OUTER JOIN templates
	ON package_templates.template_id = templates.id
LEFT OUTER JOIN package_plugins
	ON packages.id = package_plugins.package_id
LEFT OUTER JOIN plugins
	ON package_plugins.plugin_id = plugins.id`
)

func buildQuery(query, conditions string) string {
	conditions = strings.TrimSpace(conditions)
	if len(conditions) > 1 {
		query = fmt.Sprintf("%s WHERE %s", query, conditions)
	}
	return query
}

// loadAllPackages loads and returns all packages found in the database.
func (ps packageStore) queryPackage(conditions string) (*domain.Package, error) {
	query := buildQuery(defaultPackageQueryBase, conditions)
	var name, label string
	var description null.String
	err := ps.db.Raw(query).Row().Scan(&name, &label, &description)
	if err == sql.ErrNoRows {
		return nil, ErrPackageNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Package{Name: name, Label: label, Description: description.String}, nil
}

func (ps packageStore) deepQueryPackage(conditions string) (pkg *domain.Package, err error) {
	query := buildQuery(defaultPackageDeepQueryBase, conditions)
	rows, err := ps.db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			err = e
		}
	}()

	pkg = &domain.Package{}
	var gotPkgInfo bool
	for rows.Next() {
		var (
			packageName, packageLabel string
			packageDescription        null.String

			templateIsFile                                         null.Bool
			templateDestination, templatePath, templateDescription null.String

			pluginPath, pluginDescription null.String
			pluginExecNumber              null.Int
		)

		err = rows.Scan(
			&packageName,
			&packageLabel,
			&packageDescription,
			&templateIsFile,
			&templateDestination,
			&templatePath,
			&templateDescription,
			&pluginPath,
			&pluginExecNumber,
			&pluginDescription,
		)
		if err != nil {
			return nil, err
		}

		if !gotPkgInfo {
			pkg.Name = packageName
			pkg.Label = packageLabel
			pkg.Description = packageDescription.String
			gotPkgInfo = true
		}
		if templateIsFile.Valid && templateDestination.Valid && templatePath.Valid {
			pkg.Templates = append(pkg.Templates, &domain.Template{
				IsFile:      templateIsFile.Bool,
				Destination: templateDestination.String,
				Path:        templatePath.String,
				Description: templateDescription.String,
			})
		}
		if pluginPath.Valid && pluginExecNumber.Valid {
			pkg.Plugins = append(pkg.Plugins, &domain.Plugin{
				Path:        pluginPath.String,
				ExecNumber:  int(pluginExecNumber.Int64),
				Description: pluginDescription.String,
			})
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if !gotPkgInfo {
		return nil, ErrPackageNotFound
	}
	return pkg, nil
}

func (ps packageStore) queryAllLabels() ([]string, error) {
	rows, err := ps.db.Raw("SELECT label FROM packages").Rows()
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			err = e
		}
	}()
	var labels []string
	for rows.Next() {
		var label string
		err = rows.Scan(&label)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func (ps packageStore) RemovePackage(label string) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var pkg domain.Package
	tx = tx.Select("id").Where("label = ?", label).First(&pkg)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		tx.Rollback()
		return ErrPackageNotFound
	}
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	err := tx.Exec("DELETE FROM package_templates WHERE package_id = ?", pkg.ID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the actual package
	err = tx.Delete(pkg).Error
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		tx.Rollback()
		return ErrPackageNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
