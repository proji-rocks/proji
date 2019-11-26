package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
)

// sqlite represents a sqlite connection.
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
				label TEXT NOT NULL,
				is_default INTEGER NOT NULL
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
				type TEXT NOT NULL,
				run_as_sudo INTEGER NOT NULL,
				exec_num INTEGER NOT NULL,
				args TEXT
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
				is_default INTEGER NOT NULL,
				comment TEXT
			);
			INSERT INTO
				class(name, label, is_default)
			VALUES
				("unknown", "ukwn", 1);
			INSERT INTO
				project_status(title, is_default, comment)
			VALUES
				("active", 1, "Actively working on this project."),
				("inactive", 1, "Stopped working on this project for now."),
				("done", 1, "There is nothing left to do."),
				("dead", 1, "This project is dead."),
				("unknown", 1, "Status of this project is unknown.");
			CREATE UNIQUE INDEX u_class_name_idx ON class('name');
			CREATE UNIQUE INDEX u_class_label_idx ON class(label);
			CREATE UNIQUE INDEX u_class_folder_idx ON class_folder(class_id, 'target');
			CREATE UNIQUE INDEX u_class_file_idx ON class_file(class_id, 'target');
			CREATE UNIQUE INDEX u_class_script_id_name_idx ON class_script(class_id, 'name');
			CREATE UNIQUE INDEX u_class_script_type_exec_num_idx ON class_script(class_id, type, exec_num);
			CREATE UNIQUE INDEX u_project_path_idx ON project(install_path);
			CREATE UNIQUE INDEX u_status_title_idx ON project_status(title);`,
		); err != nil {
			_ = db.Close()
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

func (s *sqlite) SaveClass(class *item.Class) error {
	if err := s.saveClassInfo(class.Name, class.Label, class.IsDefault); err != nil {
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

func (s *sqlite) saveClassInfo(name, label string, isDefault bool) error {
	query := "INSERT INTO class(name, label, is_default) VALUES(?, ?, ?)"
	name = strings.ToLower(name)
	label = strings.ToLower(label)
	def := 0
	if isDefault {
		def = 1
	}
	_, err := s.db.Exec(query, name, label, def)
	return err
}

func (s *sqlite) saveFolders(classID uint, folders []*item.Folder) error {
	query := "INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, folder := range folders {
		if len(folder.Template) > 0 {
			_, err = stmt.Exec(classID, folder.Destination, folder.Template)
		} else {
			_, err = stmt.Exec(classID, folder.Destination, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveFiles(classID uint, files []*item.File) error {
	query := "INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range files {
		if len(file.Template) > 0 {
			_, err = stmt.Exec(classID, file.Destination, file.Template)
		} else {
			_, err = stmt.Exec(classID, file.Destination, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) saveScripts(classID uint, scripts []*item.Script) error {
	query := "INSERT INTO class_script(class_id, name, type, exec_num, run_as_sudo, args) VALUES(?, ?, ?, ?, ?, ?)"
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, script := range scripts {
		script.Type = strings.ToLower(script.Type)
		if script.Type != "pre" && script.Type != "post" {
			return fmt.Errorf("Script type has to be one of the following (without the single quotes): 'pre', 'post'")
		}

		args := strings.Join(script.Args, ", ")

		if script.RunAsSudo {
			_, err = stmt.Exec(classID, script.Name, script.Type, script.ExecNumber, 1, args)
		} else {
			_, err = stmt.Exec(classID, script.Name, script.Type, script.ExecNumber, 0, args)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlite) LoadClass(classID uint) (*item.Class, error) {
	name, label, isDefault, err := s.loadClassInfo(classID)
	if err != nil {
		return nil, err
	}

	class := item.NewClass(name, label, isDefault)

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

	// Assign when no errors occurred
	class.ID = classID
	class.Folders = folders
	class.Files = files
	class.Scripts = scripts
	return class, nil
}

func (s *sqlite) LoadAllClasses() ([]*item.Class, error) {
	query := "SELECT class_id FROM class ORDER BY name"

	labelRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer labelRows.Close()

	var classes []*item.Class
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

func (s *sqlite) loadClassInfo(classID uint) (string, string, bool, error) {
	query := "SELECT name, label, is_default FROM class WHERE class_id = ?"
	var name, label sql.NullString
	var isDefault sql.NullBool
	if err := s.db.QueryRow(query, classID).Scan(&name, &label, &isDefault); err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, fmt.Errorf("Class '%d' does not exist", classID)
		}
		return "", "", false, err
	}
	return name.String, label.String, isDefault.Bool, nil
}

func (s *sqlite) loadFolders(classID uint) ([]*item.Folder, error) {
	query := "SELECT target, template FROM class_folder WHERE class_id = ? ORDER BY target"

	folderRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()

	folders := make([]*item.Folder, 0)
	for folderRows.Next() {
		var dest, template sql.NullString
		if err := folderRows.Scan(&dest, &template); err != nil {
			return nil, err
		}
		folders = append(folders, &item.Folder{Destination: dest.String, Template: template.String})
	}
	return folders, nil
}

func (s *sqlite) loadFiles(classID uint) ([]*item.File, error) {
	query := "SELECT target, template FROM class_file WHERE class_id = ? ORDER BY target"

	fileRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()

	files := make([]*item.File, 0)
	for fileRows.Next() {
		var dest, template sql.NullString
		if err := fileRows.Scan(&dest, &template); err != nil {
			return nil, err
		}
		files = append(files, &item.File{Destination: dest.String, Template: template.String})
	}
	return files, nil
}

func (s *sqlite) loadScripts(classID uint) ([]*item.Script, error) {
	query := "SELECT name, type, exec_num, run_as_sudo, args FROM class_script WHERE class_id = ? ORDER BY type, exec_num"

	scriptRows, err := s.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer scriptRows.Close()

	scripts := make([]*item.Script, 0)
	for scriptRows.Next() {
		var scriptName, scriptArgs, scriptType sql.NullString
		var runAsSudo sql.NullBool
		var execNum sql.NullInt64
		if err := scriptRows.Scan(&scriptName, &scriptType, &execNum, &runAsSudo, &scriptArgs); err != nil {
			return nil, err
		}
		args := make([]string, 0)
		if scriptArgs.String != "" {
			args = strings.Split(scriptArgs.String, ", ")
		}
		scripts = append(scripts, &item.Script{
			Name:       scriptName.String,
			Type:       scriptType.String,
			ExecNumber: int(execNum.Int64),
			RunAsSudo:  runAsSudo.Bool,
			Args:       args,
		})
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

func (s *sqlite) SaveProject(proj *item.Project) error {
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

func (s *sqlite) LoadProject(projectID uint) (*item.Project, error) {
	query := "SELECT name, class_id, install_path, project_status_id FROM project WHERE project_id = ?"

	var name, installPath sql.NullString
	var classID, statusID sql.NullInt64
	if err := s.db.QueryRow(query, projectID).Scan(&name, &classID, &installPath, &statusID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Project '%d' does not exist", projectID)
		}
		return nil, err
	}

	class, err := s.LoadClass(uint(classID.Int64))
	if err != nil {
		// Load class 'unknown'
		class, err = s.LoadClass(1)
		if err != nil {
			return nil, err
		}
	}

	var status *item.Status
	status, err = s.LoadStatus(uint(statusID.Int64))
	if err != nil {
		// Load status 'unknown'
		status, err = s.LoadStatus(5)
		if err != nil {
			return nil, err
		}
	}

	return item.NewProject(projectID, name.String, installPath.String, class, status), nil
}

func (s *sqlite) LoadAllProjects() ([]*item.Project, error) {
	query := `SELECT project_id FROM project ORDER BY project_id`

	projectRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer projectRows.Close()

	var projects []*item.Project

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

func (s *sqlite) SaveStatus(status *item.Status) error {
	_, err := s.db.Exec(
		"INSERT INTO project_status(title, is_default, comment) VALUES(?, ?, ?)",
		strings.ToLower(status.Title),
		status.IsDefault,
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

func (s *sqlite) LoadStatus(statusID uint) (*item.Status, error) {
	query := "SELECT title, is_default, comment FROM project_status WHERE project_status_id = ?"
	var title, comment sql.NullString
	var isDefault sql.NullBool
	if err := s.db.QueryRow(query, statusID).Scan(&title, &isDefault, &comment); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Status '%d' does not exist", statusID)
		}
		return nil, err
	}
	return item.NewStatus(statusID, title.String, comment.String, isDefault.Bool), nil
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

func (s *sqlite) LoadAllStatuses() ([]*item.Status, error) {
	query := "SELECT project_status_id FROM project_status ORDER BY project_status_id"

	statusRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	var statuses []*item.Status

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
