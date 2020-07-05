package storage

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	mssql "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// isDatabaseDriver checks if a given database driver is actually a valid one.
func isDatabaseDriver(driver string) bool {
	var isValid bool
	switch driver {
	case sqliteDriver:
		isValid = true
	case mysqlDriver:
		isValid = true
	case mssqlDriver:
		isValid = true
	case postgresDriver:
		isValid = true
	default:
		isValid = false
	}
	return isValid
}

// Database represents a storage database. Uses gorm internally and supports sqlite, mysql, mssql and postgres
// natively.
type Database struct {
	Connection *gorm.DB
}

// newDatabaseService creates a new service instance.
func newDatabaseService(driver, connectionString string) (Service, error) {
	dialector, err := getDatabaseDialector(driver, connectionString)
	if err != nil {
		return nil, err
	}

	db := &Database{}
	db.Connection, err = gorm.Open(dialector, &gorm.Config{Logger: nil})
	db.Connection.Logger.LogMode(logger.Silent)
	return db, err
}

// getDatabaseDialector returns a sql dialector corresponding to a given driver. The dialector holds an opened
// database connection.
func getDatabaseDialector(driver, connectionString string) (gorm.Dialector, error) {
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
		dialector = nil
	}
	return dialector, err
}
