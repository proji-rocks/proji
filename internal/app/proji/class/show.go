package class

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// Show shows detailed information abour a given class
func Show(className string) error {
	className = strings.ToLower(className)

	// Connect to database
	DBDir := helper.GetConfigDir() + "/db/"
	databaseName, ok := viper.Get("database.name").(string)

	if !ok {
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

	fmt.Println(helper.ProjectHeader(className))

	if err = showLabels(tx, classID); err != nil {
		return err
	}

	if err = showFolders(tx, classID); err != nil {
		return err
	}

	if err = showFiles(tx, classID); err != nil {
		return err
	}

	if err = showScripts(tx, classID); err != nil {
		return err
	}

	return tx.Commit()
}

// showLabels shows all labels of a given class
func showLabels(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("SELECT label FROM class_label WHERE class_id = ? ORDER BY label ASC")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return err
	}
	defer query.Close()
	fmt.Println("Labels:")

	var label string
	for query.Next() {
		query.Scan(&label)
		fmt.Println(" " + label)
	}
	fmt.Println()
	return nil
}

// showFolders shows all folders of a given class
func showFolders(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return err
	}
	defer query.Close()
	fmt.Println("Folders:")

	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		fmt.Printf(" %s - %s\n", target, template)
	}
	fmt.Println()
	return nil
}

// showFiles shows all files of a given class
func showFiles(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return err
	}
	defer query.Close()
	fmt.Println("Files:")

	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		fmt.Printf(" %s - %s\n", target, template)
	}
	fmt.Println()
	return nil
}

// showScripts shows all scripts of a given class
func showScripts(tx *sql.Tx, classID int) error {
	stmt, err := tx.Prepare("SELECT name, run_as_sudo FROM class_script WHERE class_id = ? ORDER BY run_as_sudo ASC")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return err
	}
	defer query.Close()
	fmt.Println("Scripts:")

	var scriptName string
	var runAsSudo bool
	for query.Next() {
		sudo := ""
		query.Scan(&scriptName, &runAsSudo)
		if runAsSudo {
			sudo = "sudo"
		}
		fmt.Printf(" %s %s\n", scriptName, sudo)
	}
	fmt.Println()
	return nil
}
