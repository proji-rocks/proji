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
		_, err = db.Exec(
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
				install_date TEXT				
		  	);		  	
			INSERT INTO
				class(name, label, is_default)
			VALUES
				("unknown", "ukwn", 1);			
			CREATE UNIQUE INDEX u_class_name_idx ON class('name');
			CREATE UNIQUE INDEX u_class_label_idx ON class(label);
			CREATE UNIQUE INDEX u_class_folder_idx ON class_folder(class_id, 'target');
			CREATE UNIQUE INDEX u_class_file_idx ON class_file(class_id, 'target');
			CREATE UNIQUE INDEX u_class_script_id_name_idx ON class_script(class_id, 'name');
			CREATE UNIQUE INDEX u_class_script_type_exec_num_idx ON class_script(class_id, type, exec_num);
			CREATE UNIQUE INDEX u_project_path_idx ON project(install_path);`,
		)
		if err != nil {
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
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &sqlite{db, nil}, nil
}

func (s *sqlite) Close() error {
	return s.db.Close()
}

func (s *sqlite) SaveClass(class *item.Class) error {
	err := s.saveClassInfo(class.Name, class.Label, class.IsDefault)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return fmt.Errorf("class '%s' or label '%s' already exists", class.Name, class.Label)
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

	err = s.saveFolders(class.ID, class.Folders)
	if err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	err = s.saveFiles(class.ID, class.Files)
	if err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	err = s.saveScripts(class.ID, class.Scripts)
	if err != nil {
		if e := s.cancelSave(class.ID); e != nil {
			return e
		}
		return err
	}

	return s.tx.Commit()
}

func (s *sqlite) cancelSave(classID uint) error {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
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
			return fmt.Errorf("script type has to be one of the following (without the single quotes): 'pre', 'post'")
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
		err := labelRows.Scan(&classID)
		if err != nil {
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
	err := s.db.QueryRow(query, label).Scan(&classID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("class '%s' does not exist", label)
		}
		return 0, err
	}
	return uint(classID.Int64), nil
}

func (s *sqlite) loadClassInfo(classID uint) (string, string, bool, error) {
	query := "SELECT name, label, is_default FROM class WHERE class_id = ?"
	var name, label sql.NullString
	var isDefault sql.NullBool
	err := s.db.QueryRow(query, classID).Scan(&name, &label, &isDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, fmt.Errorf("class '%d' does not exist", classID)
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
		err := folderRows.Scan(&dest, &template)
		if err != nil {
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
		err := fileRows.Scan(&dest, &template)
		if err != nil {
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
		err := scriptRows.Scan(&scriptName, &scriptType, &execNum, &runAsSudo, &scriptArgs)
		if err != nil {
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
	err = s.removeScripts(classID)
	if err != nil {
		return err
	}

	err = s.removeFiles(classID)
	if err != nil {
		return err
	}

	err = s.removeFolders(classID)
	if err != nil {
		return err
	}

	err = s.removeClassInfo(classID)
	if err != nil {
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
		"INSERT INTO project(name, class_id, install_path, install_date) VALUES(?, ?, ?, ?)",
		proj.Name,
		proj.Class.ID,
		proj.InstallPath,
		t,
		1,
	)

	sqliteErr, ok := err.(sqlite3.Error)
	if ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("project already exists")
		}
	}
	return err
}

func (s *sqlite) UpdateProjectLocation(projectID uint, installPath string) error {
	_, err := s.db.Exec("UPDATE project SET install_path = ? WHERE project_id = ?", installPath, projectID)
	return err
}

func (s *sqlite) LoadProject(projectID uint) (*item.Project, error) {
	query := "SELECT name, class_id, install_path FROM project WHERE project_id = ?"

	var name, installPath sql.NullString
	var classID sql.NullInt64
	err := s.db.QueryRow(query, projectID).Scan(&name, &classID, &installPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project '%d' does not exist", projectID)
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

	return item.NewProject(projectID, name.String, installPath.String, class), nil
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

		err := projectRows.Scan(&projectID)
		if err != nil {
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
	err := s.db.QueryRow(query, installPath).Scan(&classID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("project '%s' does not exist", installPath)
		}
		return 0, err
	}
	return uint(classID.Int64), nil
}

func (s *sqlite) RemoveProject(projectID uint) error {
	_, err := s.db.Exec("DELETE FROM project WHERE project_id = ?", projectID)
	return err
}
