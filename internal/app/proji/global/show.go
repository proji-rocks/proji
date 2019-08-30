package global

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// Show shows detailed information about a global
func Show(globalType, globalID string) error {
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

	switch globalType {
	case "folder":
		err = showGlobalFolder(tx, globalID)
	case "file":
		err = showGlobalFile(tx, globalID)
	case "script":
		err = showGlobalFile(tx, globalID)
	default:
		err = fmt.Errorf("Global type not valid")
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

// showGlobalFolder shows detailed information about a global folder
func showGlobalFolder(tx *sql.Tx, globalID string) error {
	stmt, err := tx.Prepare("SELECT target_path, template_name FROM class_folder WHERE project_class_id IS NULL and class_folder_id = ? ORDER BY target_path")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(globalID)
	if err != nil {
		return err
	}
	defer query.Close()
	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		fmt.Printf(" ID: %s | Target: %s | Template: %s\n", globalID, target, template)
	}
	fmt.Println()
	return nil
}

// showGlobalFile shows detailed information about a global file
func showGlobalFile(tx *sql.Tx, globalID string) error {
	stmt, err := tx.Prepare("SELECT target_path, template_name FROM class_file WHERE project_class_id IS NULL and class_file_id = ? ORDER BY target_path")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(globalID)
	if err != nil {
		return err
	}
	defer query.Close()
	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		fmt.Printf(" ID: %s | Target: %s | Template: %s\n", globalID, target, template)
	}
	fmt.Println()
	return nil
}

// showGlobalScript shows detailed information about a global script
func showGlobalScript(tx *sql.Tx, globalID string) error {
	stmt, err := tx.Prepare("SELECT script_name, run_as_sudo FROM class_script WHERE project_class_id IS NULL and class_script_id = ? ORDER BY script_name")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(globalID)
	if err != nil {
		return err
	}
	defer query.Close()
	var scriptName string
	var runAsSudo bool
	for query.Next() {
		query.Scan(&scriptName, &runAsSudo)
		fmt.Printf(" ID: %s | Script: %s | Sudo: %v\n", globalID, scriptName, runAsSudo)
	}
	fmt.Println()
	return nil
}
