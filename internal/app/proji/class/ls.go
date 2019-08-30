package class

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// ListAll lists all classes available in the database
func ListAll() error {
	DBDir := helper.GetConfigDir() + "/db/"
	databaseName, ok := viper.Get("database.name").(string)

	if ok != true {
		return errors.New("could not read database name from config file")
	}

	db, err := sql.Open("sqlite3", DBDir+databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("SELECT class_name FROM class ORDER BY class_name ASC")
	if err != nil {
		return err
	}
	defer stmt.Close()

	queryClass, err := stmt.Query()
	if err != nil {
		return err
	}
	defer queryClass.Close()

	var className string
	for queryClass.Next() {
		queryClass.Scan(&className)
		fmt.Println(" " + className)
	}

	return tx.Commit()
}
