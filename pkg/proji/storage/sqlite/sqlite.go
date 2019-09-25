package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/nikoksr/proji/pkg/helper"

	"github.com/mattn/go-sqlite3"
	"github.com/nikoksr/proji/pkg/proji/storage"
)

// Sqlite represents a sqlite connection.
type sqlite struct {
	db *sql.DB
	tx *sql.Tx
}

// New creates a new connection to a sqlite database.
func New(path string) (storage.Service, error) {
	var db *sql.DB
	var err error

	if !helper.DoesFileExist(path) {
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}

		// Create tables
		if _, err = db.Exec(
			`CREATE TABLE IF NOT EXISTS project(
				project_id INTEGER PRIMARY KEY,
				'name' TEXT NOT NULL,
				class_id INTEGER REFERENCES class(class_id),
				install_path TEXT,
				install_date TEXT,
				project_status_id INTEGER REFERENCES project_status(project_status_id)
		  	);
		  	CREATE TABLE IF NOT EXISTS project_status(
				project_status_id INTEGER PRIMARY KEY,
				title TEXT NOT NULL,
				comment TEXT
			);
		  	CREATE TABLE IF NOT EXISTS class(
				class_id INTEGER PRIMARY KEY,
				'name' TEXT NOT NULL
		  	);
		  	CREATE TABLE IF NOT EXISTS class_label(
				class_label_id INTEGER PRIMARY KEY,
				class_id INTEGER NOT NULL REFERENCES class(class_id),
				label TEXT NOT NULL
		  	);
		  	CREATE TABLE IF NOT EXISTS class_folder(
				class_folder_id INTEGER PRIMARY KEY,
				class_id INTEGER NOT NULL REFERENCES class(class_id),
				'target' TEXT NOT NULL,
				template TEXT
		  	);
		  	CREATE TABLE IF NOT EXISTS class_file(
				class_file_id INTEGER PRIMARY KEY,
				class_id INTEGER NOT NULL REFERENCES class(class_id),
				'target' TEXT NOT NULL,
				template TEXT
		  	);
		  	CREATE TABLE IF NOT EXISTS class_script(
				class_script_id INTEGER PRIMARY KEY,
				class_id INTEGER NOT NULL REFERENCES class(class_id),
				'name' TEXT NOT NULL,
				run_as_sudo INTEGER NOT NULL
			);
			INSERT INTO
				project_status(title, comment)
			VALUES
				("active", "Actively working on this project."),
  				("inactive","Stopped working on this project for now."),
  				("done","There is nothing left to do"),
  				("dead","This project is dead.");
		  	CREATE UNIQUE INDEX u_class_idx ON class('name');
			CREATE UNIQUE INDEX u_class_label_idx ON class_label(label);
			CREATE UNIQUE INDEX u_class_folder_idx ON class_folder(class_id, 'target');
			CREATE UNIQUE INDEX u_class_file_idx ON class_file(class_id, 'target');
			CREATE UNIQUE INDEX u_class_script_idx ON class_script(class_id, 'name');
			CREATE UNIQUE INDEX u_project_path_idx ON project(install_path);`,
		); err != nil {
			db.Close()
			return nil, err
		}
	} else {
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &sqlite{db, nil}, nil
}

func (s *sqlite) Close() error {
	return s.db.Close()
}

func (s *sqlite) SaveClass(class *storage.Class) error {
	if err := s.saveName(class.Name); err != nil {
		return err
	}

	// After saving the name, the class gets a unique id.
	id, err := s.LoadClassID(class.Name)
	if err != nil {
		if e := s.cancelSave(class.Name); e != nil {
			return e
		}
		return err
	}
	class.ID = id
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx

	if err := s.saveLabels(class); err != nil {
		if e := s.cancelSave(class.Name); e != nil {
			return e
		}
		return err
	}

	if err := s.saveFolders(class); err != nil {
		if e := s.cancelSave(class.Name); e != nil {
			return e
		}
		return err
	}

	if err := s.saveFiles(class); err != nil {
		if e := s.cancelSave(class.Name); e != nil {
			return e
		}
		return err
	}

	if err := s.saveScripts(class); err != nil {
		if e := s.cancelSave(class.Name); e != nil {
			return e
		}
		return err
	}

	return s.tx.Commit()
}

func (s *sqlite) cancelSave(className string) error {
	if s.tx != nil {
		if err := s.tx.Rollback(); err != nil {
			return err
		}
	}
	return s.RemoveClass(className)
}

func (s *sqlite) saveName(name string) error {
	query := "INSERT INTO class(name) VALUES(?)"
	name = strings.ToLower(name)
	_, err := s.db.Exec(query, name)
	return err
}

func (s *sqlite) saveLabels(class *storage.Class) error {
	if s.tx == nil {
		return fmt.Errorf("No open transaction")
	}

	query := "INSERT INTO class_label(class_id, label) VALUES(?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, label := range class.Labels {
		if _, err = stmt.Exec(class.ID, strings.ToLower(label)); err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveFolders(class *storage.Class) error {
	query := "INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range class.Folders {
		if len(template) > 0 {
			_, err = stmt.Exec(class.ID, target, template)
		} else {
			_, err = stmt.Exec(class.ID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveFiles(class *storage.Class) error {
	query := "INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range class.Files {
		if len(template) > 0 {
			_, err = stmt.Exec(class.ID, target, template)
		} else {
			_, err = stmt.Exec(class.ID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveScripts(class *storage.Class) error {
	query := "INSERT INTO class_script(class_id, name, run_as_sudo) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for script, asSudo := range class.Scripts {
		if asSudo {
			_, err = stmt.Exec(class.ID, script, 1)
		} else {
			_, err = stmt.Exec(class.ID, script, 0)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) LoadClassByName(name string) (*storage.Class, error) {
	class, err := storage.NewClass(name)
	if err != nil {
		return nil, err
	}

	class.ID, err = s.LoadClassID(name)
	if err != nil {
		return nil, err
	}

	if err := s.loadLabels(class); err != nil {
		return nil, err
	}
	if err := s.loadFolders(class); err != nil {
		return nil, err
	}
	if err := s.loadFiles(class); err != nil {
		return nil, err
	}
	return class, s.loadScripts(class)
}

func (s *sqlite) LoadClassByID(id uint) (*storage.Class, error) {
	class, err := storage.NewClass("temp")
	if err != nil {
		return nil, err
	}
	class.ID = id

	if err := s.loadName(class); err != nil {
		return nil, err
	}
	if err := s.loadLabels(class); err != nil {
		return nil, err
	}
	if err := s.loadFolders(class); err != nil {
		return nil, err
	}
	if err := s.loadFiles(class); err != nil {
		return nil, err
	}
	return class, s.loadScripts(class)
}

func (s *sqlite) LoadClassID(name string) (uint, error) {
	query := "SELECT class_id FROM class WHERE name = ?"

	idRows, err := s.db.Query(query, name)
	if err != nil {
		return 0, err
	}
	defer idRows.Close()

	if !idRows.Next() {
		return 0, fmt.Errorf("Could not find class %s in storage", name)
	}

	var id uint
	err = idRows.Scan(&id)
	return id, err
}

func (s *sqlite) LoadAllClasses() ([]*storage.Class, error) {
	query := "SELECT name FROM class ORDER BY class_id"

	classRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer classRows.Close()

	var classes []*storage.Class

	for classRows.Next() {
		var name string
		classRows.Scan(&name)
		class, err := s.LoadClassByName(name)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (s *sqlite) loadName(class *storage.Class) error {
	query := "SELECT name FROM class WHERE class_id = ?"

	nameRows, err := s.db.Query(query, class.ID)
	if err != nil {
		return err
	}
	defer nameRows.Close()

	if !nameRows.Next() {
		return fmt.Errorf("Could not find class with id %d in database", class.ID)
	}
	return nameRows.Scan(&class.Name)
}

func (s *sqlite) loadLabels(class *storage.Class) error {
	query := "SELECT label FROM class_label WHERE class_id = ? ORDER BY label"

	labelRows, err := s.db.Query(query, class.ID)
	if err != nil {
		return err
	}
	defer labelRows.Close()

	for labelRows.Next() {
		var label string
		labelRows.Scan(&label)
		class.Labels = append(class.Labels, label)
	}
	return nil
}

func (s *sqlite) loadFolders(class *storage.Class) error {
	query := "SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target"

	folderRows, err := s.db.Query(query, class.ID)
	if err != nil {
		return err
	}
	defer folderRows.Close()

	for folderRows.Next() {
		var target, template string
		folderRows.Scan(&target, &template)
		class.Folders[target] = template
	}
	return nil
}

func (s *sqlite) loadFiles(class *storage.Class) error {
	query := "SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target"

	fileRows, err := s.db.Query(query, class.ID)
	if err != nil {
		return err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var target, template string
		fileRows.Scan(&target, &template)
		class.Files[target] = template
	}
	return nil
}

func (s *sqlite) loadScripts(class *storage.Class) error {
	query := "SELECT name, run_as_sudo FROM class_script WHERE class_id = ? ORDER BY run_as_sudo, name"

	scriptRows, err := s.db.Query(query, class.ID)
	if err != nil {
		return err
	}
	defer scriptRows.Close()

	for scriptRows.Next() {
		var scriptName string
		var runAsSudo bool
		scriptRows.Scan(&scriptName, &runAsSudo)
		class.Scripts[scriptName] = runAsSudo
	}
	return nil
}

func (s *sqlite) RemoveClass(name string) error {
	var err error

	classID, err := s.LoadClassID(name)
	if err != nil {
		return err
	}

	s.tx, err = s.db.Begin()
	if err != nil {
		return err
	}

	// Remove class and dependencies
	if err = s.removeName(classID); err != nil {
		return err
	}
	if err = s.removeLabels(classID); err != nil {
		return err
	}
	if err = s.removeFolders(classID); err != nil {
		return err
	}
	if err = s.removeFiles(classID); err != nil {
		return err
	}
	if err = s.removeScripts(classID); err != nil {
		return err
	}
	return s.tx.Commit()
}

func (s *sqlite) removeName(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class WHERE class_id = ?", classID)
	return err
}

func (s *sqlite) removeLabels(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class_label WHERE class_id = ?", classID)
	return err
}

func (s *sqlite) removeFolders(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class_folder WHERE class_id = ?", classID)
	return err
}

func (s *sqlite) removeFiles(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class_file WHERE class_id = ?", classID)
	return err
}

func (s *sqlite) removeScripts(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class_script WHERE class_id = ?", classID)
	return err
}

func (s *sqlite) DoesLabelExist(label string) (uint, error) {
	query := "SELECT class_id FROM class_label WHERE label = ?"
	var id uint
	err := s.db.QueryRow(query, label).Scan(&id)
	return id, err
}

func (s *sqlite) TrackProject(proj *storage.Project) error {
	t := time.Now().Local()
	_, err := s.db.Exec(
		"INSERT INTO project(name, class_id, install_path, install_date, project_status_id) VALUES(?, ?, ?, ?, ?)",
		proj.Name,
		proj.Class.ID,
		proj.InstallPath,
		t,
		1,
	)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("Project already exists")
		}
	}
	return err
}

func (s *sqlite) UntrackProject(projectID uint) error {
	_, err := s.db.Exec("DELETE FROM project WHERE project_id = ?", projectID)
	return err
}

func (s *sqlite) UpdateProjectStatus(projectID, statusID uint) error {
	_, err := s.db.Exec("UPDATE project SET project_status_id = ? WHERE project_id = ?", statusID, projectID)
	return err
}

func (s *sqlite) UpdateProjectLocation(projectID uint, installPath string) error {
	_, err := s.db.Exec("UPDATE project SET install_path = ? WHERE project_id = ?", installPath, projectID)
	return err
}

func (s *sqlite) LoadProjectID(path string) (uint, error) {
	query := "SELECT project_id FROM project WHERE install_path = ?"

	idRows, err := s.db.Query(query, path)
	if err != nil {
		return 0, err
	}
	defer idRows.Close()

	if !idRows.Next() {
		return 0, fmt.Errorf("Could not find project %s in database", path)
	}

	var id uint
	err = idRows.Scan(&id)
	return id, err
}

func (s *sqlite) ListProjects() ([]*storage.Project, error) {
	query := `
		SELECT
			p.project_id,
			p.name as title,
			p.install_path,
			c.name as class_name,
			ps.title
		FROM
			project p
		JOIN project_status ps ON
			p.project_status_id = ps.project_status_id
		JOIN class c ON
			p.class_id = c.class_id
		ORDER BY
			p.project_id
	`

	projectRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer projectRows.Close()

	var projects []*storage.Project

	for projectRows.Next() {
		var project storage.Project
		var status storage.Status
		var class storage.Class

		if err := projectRows.Scan(&project.ID, &project.Name, &project.InstallPath, &class.Name, &status.Title); err != nil {
			return nil, err
		}
		project.Status = &status
		project.Class = &class
		projects = append(projects, &project)
	}
	return projects, nil
}

func (s *sqlite) AddProjectStatus(status *storage.Status) error {
	_, err := s.db.Exec(
		"INSERT INTO project(title, comment) VALUES(?, ?)",
		status.Title,
		status.Comment,
	)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("Status already exists")
		}
	}
	return err
}

func (s *sqlite) RemoveProjectStatus(statusID uint) error {
	_, err := s.db.Exec("DELETE FROM project_status WHERE project_status_id = ?", statusID)
	return err
}

func (s *sqlite) ListAvailableProjectStatuses() ([]*storage.Status, error) {
	query := "SELECT * FROM project_status ORDER BY project_status_id"

	statusRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	var statuses []*storage.Status

	for statusRows.Next() {
		var status storage.Status
		if err := statusRows.Scan(&status.ID, &status.Title, &status.Comment); err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	return statuses, nil
}
