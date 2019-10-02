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

	if !helper.DoesPathExist(path) {
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}

		// Create tables
		if _, err = db.Exec(
			`CREATE TABLE IF NOT EXISTS class(
				class_id INTEGER PRIMARY KEY,
				'name' TEXT NOT NULL,
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
			CREATE TABLE IF NOT EXISTS project(
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
			INSERT INTO
				project_status(title, comment)
			VALUES
				("active", "Actively working on this project."),
  				("inactive", "Stopped working on this project for now."),
  				("done", "There is nothing left to do"),
				("dead", "This project is dead."),
				("unknown", "Status of this project is unknown.");
			CREATE UNIQUE INDEX u_class_name_idx ON class('name');
			CREATE UNIQUE INDEX u_class_label_idx ON class(label);
			CREATE UNIQUE INDEX u_class_folder_idx ON class_folder(class_id, 'target');
			CREATE UNIQUE INDEX u_class_file_idx ON class_file(class_id, 'target');
			CREATE UNIQUE INDEX u_class_script_idx ON class_script(class_id, 'name');
			CREATE UNIQUE INDEX u_project_path_idx ON project(install_path);
			CREATE UNIQUE INDEX u_status_title_idx ON project_status(title);`,
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
	if err := s.saveClassInfo(class.Name, class.Label); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return fmt.Errorf("Class '%s' or label '%s' already exists", class.Name, class.Label)
			}
		}
		return err
	}

	// After saving the class info the class gets a unique ID, load it
	id, err := s.LoadClassIDByLabel(class.Label)
	if err != nil {
		return err
	}
	class.ID = id

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx

	if err := s.saveFolders(class.ID, class.Folders); err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	if err := s.saveFiles(class.ID, class.Files); err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	if err := s.saveScripts(class.ID, class.Scripts); err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	return s.tx.Commit()
}

func (s *sqlite) cancelSave(classID uint) error {
	if s.tx != nil {
		if err := s.tx.Rollback(); err != nil {
			return err
		}
	}
	return s.RemoveClass(classID)
}

func (s *sqlite) saveClassInfo(name, label string) error {
	query := "INSERT INTO class(name, label) VALUES(?, ?)"
	name = strings.ToLower(name)
	label = strings.ToLower(label)
	_, err := s.db.Exec(query, name, label)
	return err
}

func (s *sqlite) saveFolders(classID uint, folders map[string]string) error {
	query := "INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range folders {
		if len(template) > 0 {
			_, err = stmt.Exec(classID, target, template)
		} else {
			_, err = stmt.Exec(classID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveFiles(classID uint, files map[string]string) error {
	query := "INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, template := range files {
		if len(template) > 0 {
			_, err = stmt.Exec(classID, target, template)
		} else {
			_, err = stmt.Exec(classID, target, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveScripts(classID uint, scripts map[string]bool) error {
	query := "INSERT INTO class_script(class_id, name, run_as_sudo) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for script, asSudo := range scripts {
		if asSudo {
			_, err = stmt.Exec(classID, script, 1)
		} else {
			_, err = stmt.Exec(classID, script, 0)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) LoadClass(classID uint) (*storage.Class, error) {
	name, label, err := s.loadClassInfo(classID)
	if err != nil {
		return nil, err
	}

	class, err := storage.NewClass(name, label)
	if err != nil {
		return nil, err
	}

	folders, err := s.loadFolders(classID)
	if err != nil {
		return nil, err
	}
	files, err := s.loadFiles(classID)
	if err != nil {
		return nil, err
	}
	scripts, err := s.loadScripts(classID)
	if err != nil {
		return nil, err
	}

	// Assign when no errors occured
	class.ID = classID
	class.Folders = folders
	class.Files = files
	class.Scripts = scripts
	return class, nil
}

func (s *sqlite) LoadAllClasses() ([]*storage.Class, error) {
	query := "SELECT class_id FROM class ORDER BY name"

	labelRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer labelRows.Close()

	var classes []*storage.Class
	for labelRows.Next() {
		var classID sql.NullInt64
		if err := labelRows.Scan(&classID); err != nil {
			return nil, err
		}
		class, err := s.LoadClass(uint(classID.Int64))
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (s *sqlite) LoadClassIDByLabel(label string) (uint, error) {
	query := "SELECT class_id FROM class WHERE label = ?"
	var classID sql.NullInt64
	if err := s.db.QueryRow(query, label).Scan(&classID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("Class '%s' does not exist", label)
		}
		return 0, err
	}
	return uint(classID.Int64), nil
}

func (s *sqlite) loadClassInfo(classID uint) (string, string, error) {
	query := "SELECT name, label FROM class WHERE class_id = ?"
	var name, label sql.NullString
	if err := s.db.QueryRow(query, classID).Scan(&name, &label); err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("Class '%d' does not exist", classID)
		}
		return "", "", err
	}
	return name.String, label.String, nil
}

func (s *sqlite) loadFolders(classID uint) (map[string]string, error) {
	query := "SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target"

	folderRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()

	folders := make(map[string]string)
	for folderRows.Next() {
		var target, template sql.NullString
		if err := folderRows.Scan(&target, &template); err != nil {
			return nil, err
		}
		folders[target.String] = template.String
	}
	return folders, nil
}

func (s *sqlite) loadFiles(classID uint) (map[string]string, error) {
	query := "SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target"

	fileRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()

	files := make(map[string]string)
	for fileRows.Next() {
		var target, template sql.NullString
		if err := fileRows.Scan(&target, &template); err != nil {
			return nil, err
		}
		files[target.String] = template.String
	}
	return files, nil
}

func (s *sqlite) loadScripts(classID uint) (map[string]bool, error) {
	query := "SELECT name, run_as_sudo FROM class_script WHERE class_id = ? ORDER BY class_script_id"

	scriptRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer scriptRows.Close()

	scripts := make(map[string]bool)
	for scriptRows.Next() {
		var scriptName sql.NullString
		var runAsSudo sql.NullBool
		if err := scriptRows.Scan(&scriptName, &runAsSudo); err != nil {
			return nil, err
		}
		scripts[scriptName.String] = runAsSudo.Bool
	}
	return scripts, nil
}

func (s *sqlite) RemoveClass(classID uint) error {
	var err error
	s.tx, err = s.db.Begin()
	if err != nil {
		return err
	}

	// Remove class and dependencies
	if err = s.removeScripts(classID); err != nil {
		return err
	}
	if err = s.removeFiles(classID); err != nil {
		return err
	}
	if err = s.removeFolders(classID); err != nil {
		return err
	}
	if err = s.removeClassInfo(classID); err != nil {
		return err
	}
	return s.tx.Commit()
}

func (s *sqlite) removeClassInfo(classID uint) error {
	_, err := s.tx.Exec("DELETE FROM class WHERE class_id = ?", classID)
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

func (s *sqlite) SaveProject(proj *storage.Project) error {
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

func (s *sqlite) UpdateProjectStatus(projectID, statusID uint) error {
	_, err := s.db.Exec("UPDATE project SET project_status_id = ? WHERE project_id = ?", statusID, projectID)
	return err
}

func (s *sqlite) UpdateProjectLocation(projectID uint, installPath string) error {
	_, err := s.db.Exec("UPDATE project SET install_path = ? WHERE project_id = ?", installPath, projectID)
	return err
}

func (s *sqlite) LoadProject(projectID uint) (*storage.Project, error) {
	query := "SELECT name, class_id, install_path, project_status_id FROM project WHERE project_id = ?"

	var name, installPath sql.NullString
	var classID, statusID sql.NullInt64
	if err := s.db.QueryRow(query, projectID).Scan(&name, &classID, &installPath, &statusID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Project '%d' does not exist", projectID)
		}
		return nil, err
	}
	project, err := storage.NewProject(projectID, name.String, installPath.String, uint(classID.Int64), uint(statusID.Int64), s)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *sqlite) LoadAllProjects() ([]*storage.Project, error) {
	query := `SELECT project_id FROM project ORDER BY project_id`

	projectRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer projectRows.Close()

	var projects []*storage.Project

	for projectRows.Next() {
		var projectID sql.NullInt64

		if err := projectRows.Scan(&projectID); err != nil {
			return nil, err
		}

		project, err := s.LoadProject(uint(projectID.Int64))
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (s *sqlite) LoadProjectID(installPath string) (uint, error) {
	query := "SELECT project_id FROM project WHERE install_path = ?"
	var classID sql.NullInt64
	if err := s.db.QueryRow(query, installPath).Scan(&classID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("Project '%s' does not exist", installPath)
		}
		return 0, err
	}
	return uint(classID.Int64), nil
}

func (s *sqlite) RemoveProject(projectID uint) error {
	_, err := s.db.Exec("DELETE FROM project WHERE project_id = ?", projectID)
	return err
}

func (s *sqlite) SaveStatus(status *storage.Status) error {
	_, err := s.db.Exec(
		"INSERT INTO project_status(title, comment) VALUES(?, ?)",
		strings.ToLower(status.Title),
		status.Comment,
	)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("Status '%s' already exists", status.Title)
		}
	}
	return err
}

func (s *sqlite) UpdateStatus(statusID uint, title, comment string) error {
	_, err := s.db.Exec("UPDATE project_status SET title = ?, comment = ? WHERE project_status_id = ?",
		strings.ToLower(title),
		comment,
		statusID,
	)
	return err
}

func (s *sqlite) LoadStatus(statusID uint) (*storage.Status, error) {
	query := "SELECT title, comment FROM project_status WHERE project_status_id = ?"
	var title, comment sql.NullString
	if err := s.db.QueryRow(query, statusID).Scan(&title, &comment); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Status '%d' does not exist", statusID)
		}
		return nil, err
	}
	return storage.NewStatus(statusID, title.String, comment.String), nil
}

func (s *sqlite) LoadStatusID(title string) (uint, error) {
	query := "SELECT project_status_id FROM project_status WHERE title = ?"
	var statusID sql.NullInt64
	if err := s.db.QueryRow(query, title).Scan(&statusID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("Status '%s' does not exist", title)
		}
		return 0, err
	}
	return uint(statusID.Int64), nil
}

func (s *sqlite) LoadAllStatuses() ([]*storage.Status, error) {
	query := "SELECT project_status_id FROM project_status ORDER BY project_status_id"

	statusRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	var statuses []*storage.Status

	for statusRows.Next() {
		var statusID sql.NullInt64
		if err = statusRows.Scan(&statusID); err != nil {
			return nil, err
		}
		status, err := s.LoadStatus(uint(statusID.Int64))
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func (s *sqlite) RemoveStatus(statusID uint) error {
	_, err := s.db.Exec("DELETE FROM project_status WHERE project_status_id = ?", statusID)
	return err
}
