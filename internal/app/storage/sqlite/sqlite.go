package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nikoksr/proji/internal/app/proji/class"
	"github.com/nikoksr/proji/internal/app/storage"
)

// Sqlite represents a sqlite connection.
type sqlite struct {
	DB *sql.DB
	Tx *sql.Tx
}

// New creates a new connection to a sqlite database.
func New(name string) (storage.Service, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &sqlite{db, nil}, nil
}

// Close the sqlite connection.
func (s *sqlite) Close() error {
	return s.Close()
}

func (s *sqlite) Save(c *class.Class) error {
	if err := s.saveName(c.Name); err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	s.Tx = tx

	if err := s.saveLabels(c); err != nil {
		return err
	}

	if err := s.saveFolders(c); err != nil {
		return err
	}

	if err := s.saveFiles(c); err != nil {
		return err
	}

	if err := s.saveScripts(c); err != nil {
		return err
	}

	return s.Tx.Commit()
}

func (s *sqlite) saveName(name string) error {
	query := "INSERT INTO class(name) VALUES(?)"
	name = strings.ToLower(name)
	_, err := s.DB.Exec(query, name)
	return err
}

func (s *sqlite) saveLabels(c *class.Class) error {
	if s.Tx == nil {
		return fmt.Errorf("no open transaction")
	}

	query := "INSERT INTO class_label(class_id, label) VALUES(?, ?)"
	stmt, err := s.Tx.Prepare(query)
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

func (s *sqlite) saveFolders(c *class.Class) error {
	query := "INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.Tx.Prepare(query)
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

func (s *sqlite) saveFiles(c *class.Class) error {
	query := "INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.Tx.Prepare(query)
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

func (s *sqlite) saveScripts(c *class.Class) error {
	query := "INSERT INTO class_script(class_id, name, run_as_sudo) VALUES(?, ?, ?)"
	stmt, err := s.Tx.Prepare(query)
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
func (s *sqlite) Load(name string) (Item, error) {
	c := class.New(name)

	if err := s.loadID(c); err != nil {
		return nil, err
	}
	if err := s.loadLabels(c); err != nil {
		return nil, err
	}
	if err := s.loadFolders(c); err != nil {
		return nil, err
	}
	if err := s.loadFiles(c); err != nil {
		return nil, err
	}
	return c, s.loadScripts(c)
}

// LoadID loads the class id.
func (s *sqlite) LoadID(c *class.Class) error {
	query := "SELECT class_id FROM class WHERE name = ?"

	idRows, err := s.DB.Query(query, c.Name)
	if err != nil {
		return err
	}
	defer id.Close()

	if !resID.Next() {
		return fmt.Errorf("could not find class %s in database", c.Name)
	}
	return resID.Scan(&c.ID)
}

// loadLabels loads all labels of a given class.
func (s *sqlite) loadLabels(c *class.Class) error {
	query := "SELECT label FROM class_label WHERE class_id = ? ORDER BY label"

	labelRows, err := s.DB.Query(query, c.ID)
	if err != nil {
		return err
	}
	defer labelRows.Close()

	for labelRows.Next() {
		var label string
		labelRows.Scan(&label)
		c.Labels = append(c.Labels, label)
	}
	return nil
}

// loadFolders loads all folders of a given class.
func (s *sqlite) loadFolders(c *class.Class) error {
	query := "SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target"

	folderRows, err := s.DB.Query(query, c.ID)
	if err != nil {
		return err
	}
	defer folderRows.Close()

	for folderRows.Next() {
		var target, template string
		folderRows.Scan(&target, &template)
		c.Folders[target] = template
	}
	return nil
}

// loadFiles loads all files of a given class
func (s *sqlite) loadFiles(c *class.Class) error {
	query := "SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target"

	fileRows, err := s.DB.Query(query, c.ID)
	if err != nil {
		return err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var target, template string
		fileRows.Scan(&target, &template)
		c.Files[target] = template
	}
	return nil
}

// loadScripts loads all scripts of a given class
func (s *sqlite) loadScripts(c *class.Class) error {
	query := "SELECT name, run_as_sudo FROM class_script WHERE class_id = ? ORDER BY run_as_sudo, name"

	scriptRows, err := s.DB.Query(query, c.ID)
	if err != nil {
		return err
	}
	defer scriptRows.Close()

	for scriptRows.Next() {
		var scriptName string
		var runAsSudo bool
		scriptRows.Scan(&scriptName, &runAsSudo)
		c.Scripts[scriptName] = runAsSudo
	}
	return nil
}
