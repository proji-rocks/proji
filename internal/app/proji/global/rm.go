package global

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// RemoveGlobal removes a global from the database
func RemoveGlobal(globalType string, globalID []string) error {
	if len(globalID) < 1 {
		return fmt.Errorf("no global id given")
	}

	// Connect to database
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

	// Define type specific statement
	var stmt *sql.Stmt
	switch globalType {
	case "folder":
		stmt, err = tx.Prepare("DELETE FROM class_folder WHERE class_folder_id = ? AND class_id IS NULL")
		if err != nil {
			return err
		}
	case "file":
		stmt, err = tx.Prepare("DELETE FROM class_file WHERE class_file_id = ? AND class_id IS NULL")
		if err != nil {
			return err
		}
	case "script":
		stmt, err = tx.Prepare("DELETE FROM class_script WHERE class_script_id = ? AND class_id IS NULL")
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("global type not valid")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Remove globals
	for _, glob := range globalID {
		id, err := strconv.Atoi(glob)
		if err != nil {
			fmt.Printf("> %s is not an id", glob)
			continue
		}
		if _, err = stmt.Exec(id); err != nil {
			fmt.Printf("> Removing %d: %e.\n", id, err)
		}
	}

	err = tx.Commit()
	return err
}
