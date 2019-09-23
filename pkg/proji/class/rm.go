package class

import (
	"database/sql"
	"errors"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/viper"
)

// Remove removes an existing class and all of its depending settings in other tables from the database.
func (c *Class) Remove() error {

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

	if err = c.loadID(db); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Remove class and dependencies
	if err = c.removeName(tx); err != nil {
		return err
	}
	if err = c.removeLabels(tx); err != nil {
		return err
	}
	if err = c.removeFolders(tx); err != nil {
		return err
	}
	if err = c.removeFiles(tx); err != nil {
		return err
	}
	if err = c.removeScripts(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// removeName removes an existing class name.
func (c *Class) removeName(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM class WHERE class_id = ?", c.ID)
	return err
}

// removeLabels removes all class labels.
func (c *Class) removeLabels(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM class_label WHERE class_id = ?", c.ID)
	return err
}

// removeFolders removes all class folders.
func (c *Class) removeFolders(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM class_folder WHERE class_id = ?", c.ID)
	return err
}

// removeFiles removes all class files.
func (c *Class) removeFiles(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM class_file WHERE class_id = ?", c.ID)
	return err
}

// removeScripts removes all class scripts.
func (c *Class) removeScripts(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM class_script WHERE class_id = ?", c.ID)
	return err
}
