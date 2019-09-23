package class

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// Class struct represents a proji class
type Class struct {
	// The class name
	Name string

	// The class ID
	ID uint

	// All class related labels
	Labels []string

	// All class related folders
	Folders map[string]string

	// All class related files
	Files map[string]string

	// All class related scripts
	Scripts map[string]bool
}

// loadID loads the id of a given class name
func (c *Class) loadID(db *sql.DB) error {
	resID, err := db.Query("SELECT class_id FROM class WHERE name = ?", c.Name)
	if err != nil {
		return err
	}
	defer resID.Close()

	if !resID.Next() {
		return fmt.Errorf("could not find class %s in database", c.Name)
	}
	return resID.Scan(&c.ID)
}

// New returns a new class
func New(name string) *Class {
	return &Class{
		Name:    name,
		ID:      0,
		Labels:  make([]string, 0),
		Folders: make(map[string]string),
		Files:   make(map[string]string),
		Scripts: make(map[string]bool),
	}
}

// Save saves a new class and its data to the database.
func (c *Class) Save() error {
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

	if err = c.saveName(db); err != nil {
		return err
	}

	// Necessary for save operations
	if err = c.loadID(db); err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err = c.saveLabels(tx); err != nil {
		return err
	}

	if err = c.saveFolders(tx); err != nil {
		return err
	}

	if err = c.saveFiles(tx); err != nil {
		return err
	}

	if err = c.saveScripts(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// saveName saves a new class name to the database and gets the associated class id
func (c *Class) saveName(db *sql.DB) error {
	c.Name = strings.ToLower(c.Name)
	_, err := db.Exec("INSERT INTO class(name) VALUES(?)", c.Name)
	return err
}

// saveLabels saves new class labels to the database
func (c *Class) saveLabels(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO class_label(class_id, label) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, label := range c.Labels {
		if _, err = stmt.Exec(c.ID, strings.ToLower(label)); err != nil {
			return err
		}
	}
	return nil
}

// saveFolders saves new class folders to the database
func (c *Class) saveFolders(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range c.Folders {
		if len(template) > 0 {
			_, err = stmt.Exec(c.ID, target, template)
		} else {
			_, err = stmt.Exec(c.ID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// saveFiles saves new class files to the database
func (c *Class) saveFiles(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range c.Files {
		if len(template) > 0 {
			_, err = stmt.Exec(c.ID, target, template)
		} else {
			_, err = stmt.Exec(c.ID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// saveScripts saves new class scripts to the database
func (c *Class) saveScripts(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO class_script(class_id, name, run_as_sudo) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for script, asSudo := range c.Scripts {
		if asSudo {
			_, err = stmt.Exec(c.ID, script, 1)
		} else {
			_, err = stmt.Exec(c.ID, script, 0)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads a class struct from the database if given a valid name.
func (c *Class) Load() error {
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
	if err = c.loadLabels(db); err != nil {
		return err
	}
	if err = c.loadFolders(db); err != nil {
		return err
	}
	if err = c.loadFiles(db); err != nil {
		return err
	}
	return c.loadScripts(db)
}

// loadLabels loads all labels of a given class.
func (c *Class) loadLabels(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT label FROM class_label WHERE class_id = ? ORDER BY label")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(c.ID)
	if err != nil {
		return err
	}
	defer query.Close()

	for query.Next() {
		var label string
		query.Scan(&label)
		c.Labels = append(c.Labels, label)
	}
	return nil
}

// loadFolders loads all folders of a given class.
func (c *Class) loadFolders(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(c.ID)
	if err != nil {
		return err
	}
	defer query.Close()

	for query.Next() {
		var target, template string
		query.Scan(&target, &template)
		c.Folders[target] = template
	}
	return nil
}

// loadFiles loads all files of a given class
func (c *Class) loadFiles(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(c.ID)
	if err != nil {
		return err
	}
	defer query.Close()

	for query.Next() {
		var target, template string
		query.Scan(&target, &template)
		c.Files[target] = template
	}
	return nil
}

// loadScripts loads all scripts of a given class
func (c *Class) loadScripts(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT name, run_as_sudo FROM class_script WHERE class_id = ? ORDER BY run_as_sudo, name")
	if err != nil {
		return err
	}
	defer stmt.Close()

	query, err := stmt.Query(c.ID)
	if err != nil {
		return err
	}
	defer query.Close()

	for query.Next() {
		var scriptName string
		var runAsSudo bool
		query.Scan(&scriptName, &runAsSudo)
		c.Scripts[scriptName] = runAsSudo
	}
	return nil
}
