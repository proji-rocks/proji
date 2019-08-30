package global

import (
	"database/sql"
	"fmt"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// ListAll lists all proji globals
func ListAll() error {
	DBDir := helper.GetConfigDir() + "/db/"
	databaseName, ok := viper.Get("database.name").(string)

	if ok != true {
		return fmt.Errorf("could not read database name from config file")
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

	if err = showGlobalFolders(tx); err != nil {
		return nil
	}

	if err = showGlobalFiles(tx); err != nil {
		return nil
	}

	if err = showGlobalScripts(tx); err != nil {
		return nil
	}

	return tx.Commit()
}

// showGlobalFolders shows all global folders
func showGlobalFolders(tx *sql.Tx) error {
	stmt, err := tx.Prepare("SELECT class_folder_id, target_path, template_name FROM class_folder WHERE class_id is NULL ORDER BY target_path")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query()
	if err != nil {
		return err
	}
	defer query.Close()

	fmt.Println("Folders:")
	for query.Next() {
		var id, target, template string
		query.Scan(&id, &target, &template)
		fmt.Printf(" ID: %s | Target: %s | Template: %s\n", id, target, template)
	}
	fmt.Println()
	return nil
}

// showGlobalFiles shows all global files
func showGlobalFiles(tx *sql.Tx) error {
	stmt, err := tx.Prepare("SELECT class_file_id, target_path, template_name FROM class_file WHERE class_id is NULL ORDER BY target_path")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query()
	if err != nil {
		return err
	}
	defer query.Close()

	fmt.Println("Files:")
	for query.Next() {
		var id, target, template string
		query.Scan(&id, &target, &template)
		fmt.Printf(" ID: %s | Target: %s | Template: %s\n", id, target, template)
	}
	fmt.Println()
	return nil
}

// showGlobalScripts shows all global scripts
func showGlobalScripts(tx *sql.Tx) error {
	stmt, err := tx.Prepare("SELECT class_script_id, script_name, run_as_sudo FROM class_script WHERE class_id is NULL ORDER BY script_name")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query()
	if err != nil {
		return err
	}
	defer query.Close()

	fmt.Println("Scripts:")
	for query.Next() {
		var id, scriptName string
		var runAsSudo bool
		query.Scan(&id, &scriptName, &runAsSudo)
		fmt.Printf(" ID: %s | Script: %s | Sudo: %v\n", id, scriptName, runAsSudo)
	}
	fmt.Println()
	return nil
}
