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
		return err
	}

	if err = showGlobalFiles(tx); err != nil {
		return err
	}

	if err = showGlobalScripts(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// showGlobalFolders shows all global folders
func showGlobalFolders(tx *sql.Tx) error {
	stmt, err := tx.Prepare("SELECT global_folder_id, target, template FROM global_folder ORDER BY global_folder_id ASC")
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
	stmt, err := tx.Prepare("SELECT global_file_id, target, template FROM global_file ORDER BY global_file_id ASC")
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
	stmt, err := tx.Prepare("SELECT global_script_id, name, run_as_sudo FROM global_script ORDER BY global_script_id")
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
