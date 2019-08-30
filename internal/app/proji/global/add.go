package global

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// AddGlobal adds a new global to the proji database
func AddGlobal(globalType string, newGlobal []string) error {

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

	// Insert data
	tx, err := db.Begin()

	switch globalType {
	case "folder":
		err = insertGlobalFolder(tx, newGlobal)
	case "file":
		err = insertGlobalFile(tx, newGlobal)
	case "script":
		err = insertGlobalFile(tx, newGlobal)
	default:
		err = fmt.Errorf("Global type not valid")
	}

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("> Added global %s successfully.\n", newGlobal)
	return nil
}

// insertGlobalFolder inserts a new global folder into the database
func insertGlobalFolder(tx *sql.Tx, folder []string) error {
	stmt, err := tx.Prepare("INSERT INTO global_folder(target, template) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	target := folder[0]

	if len(folder) > 1 {
		_, err = stmt.Exec(target, folder[1])
	} else {
		_, err = stmt.Exec(target, nil)
	}
	return err
}

// insertGlobalFile inserts a new global file into the database
func insertGlobalFile(tx *sql.Tx, file []string) error {
	stmt, err := tx.Prepare("INSERT INTO global_file(target, template) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	target := file[0]

	if len(file) > 1 {
		_, err = stmt.Exec(target, file[1])
	} else {
		_, err = stmt.Exec(target, nil)
	}
	return err
}

// insertGlobalScript inserts a new global script into the database
func insertGlobalScript(tx *sql.Tx, script []string) error {
	stmt, err := tx.Prepare("INSERT INTO global_script(name, run_as_sudo) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(script[0], script[1])
	return err
}
