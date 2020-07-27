package database

import (
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	mssql "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	sqliteDriver   = "sqlite3"
	mysqlDriver    = "mysql"
	mssqlDriver    = "mssql"
	postgresDriver = "postgres"
)

// Database represents a storage database. Uses gorm internally and supports sqlite, mysql, mssql and postgres
// natively.
type Database struct {
	Connection *gorm.DB
}

// New creates a new database instance.
func New(driver, connectionString string) (*Database, error) {
	// Try to get sql dialect from driver; if successful connect to db with Connection string
	dialector, err := getDialector(driver, connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "get sql dialect")
	}

	// Create database struct; gorm's logger turned off for now by default
	db := &Database{}
	db.Connection, err = gorm.Open(
		dialector,
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger: logger.New(
				nil,
				logger.Config{
					Colorful: false,
					LogLevel: logger.Silent,
				},
			),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "open Connection")
	}
	return db, nil
}

func (db Database) Migrate() error {
	modelList := []interface{}{
		&domain.Package{},
		&domain.Plugin{},
		&domain.Project{},
		&domain.Template{},
	}
	for _, model := range modelList {
		err := db.Connection.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to auto-migrate domain, %s", err.Error())
		}
	}
	return nil
}

// getDialector returns a sql dialector corresponding to a given driver. The dialector holds an opened
// database Connection.
func getDialector(driver, connectionString string) (gorm.Dialector, error) {
	var err error
	var dialector gorm.Dialector
	switch driver {
	case sqliteDriver:
		dialector = sqlite.Open(connectionString)
	case mysqlDriver:
		dialector = mysql.Open(connectionString)
	case mssqlDriver:
		dialector = mssql.Open(connectionString)
	case postgresDriver:
		dialector = postgres.Open(connectionString)
	default:
		err = &UnsupportedDatabaseDialectError{Dialect: driver}
	}
	return dialector, err
}

// UnsupportedDatabaseDialectError represents an error for the case that the user passed a db driver for a
// unsupported database dialect. Proji uses Gorm under the hood so take a look at its docs for a list of
// documented dialects.
// https://gorm.io/docs/connecting_to_the_database.html#Supported-Databases
type UnsupportedDatabaseDialectError struct {
	Dialect string
}

func (e *UnsupportedDatabaseDialectError) Error() string {
	return fmt.Sprintf("%s is not in the list of supported database dialects", e.Dialect)
}
