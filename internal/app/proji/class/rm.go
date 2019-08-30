package class

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// RemoveClass removes an existing class and all of its depending settings in other tables from the database
func RemoveClass(className string) error {
	className = strings.ToLower(className)

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

	// Get class id
	classID, err := helper.QueryClassID(tx, className)
	if err != nil {
		return err
	}

	// Remove class and dependencies
	err = removeClass(tx, classID)
	if err != nil {
		return err
	}
	err = removeLabels(tx, classID)
	if err != nil {
		return err
	}
	err = removeFolders(tx, classID)
	if err != nil {
		return err
	}
	err = removeFiles(tx, classID)
	if err != nil {
		return err
	}
	err = removeScripts(tx, classID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("> Removed class %s successfully.\n", className)
	return nil
}

// removeClass removes an existing class from the database
func removeClass(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("DELETE FROM class WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classID)
	if err != nil {
		return err
	}
	return nil
}

// removeLabels removes all class labels from the database
func removeLabels(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("DELETE FROM class_label WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classID)
	if err != nil {
		return err
	}
	return nil
}

// removeFolders removes all class folders from the database
func removeFolders(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("DELETE FROM class_folder WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classID)
	if err != nil {
		return err
	}
	return nil
}

// removeFiles removes all class files from the database
func removeFiles(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("DELETE FROM class_file WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classID)
	if err != nil {
		return err
	}
	return nil
}

// removeScripts removes all class scripts from the database
func removeScripts(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("DELETE FROM class_script WHERE class_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classID)
	if err != nil {
		return err
	}
	return nil
}
