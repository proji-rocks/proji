package class

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// Export exports a given class to a toml config file
func Export(className string) error {
	className = strings.ToLower(className)

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

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Get class id
	classID, err := helper.QueryClassID(tx, className)
	if err != nil {
		return err
	}

	// Get data to export
	title := "proji class " + className

	labels, err := exportLabels(tx, classID)
	if err != nil {
		return nil
	}

	folders, err := exportFolders(tx, classID)
	if err != nil {
		return nil
	}

	files, err := exportFiles(tx, classID)
	if err != nil {
		return nil
	}

	scripts, err := exportScripts(tx, classID)
	if err != nil {
		return nil
	}

	// Create config string
	var configTxt = map[string]interface{}{
		"title": title,
		"class": map[string]string{
			"name": className,
		},
		"labels": map[string][]string{
			"data": labels,
		},
		"folders": folders,
		"files":   files,
		"scripts": scripts,
	}

	// Export data to toml
	confName := "proji-export-" + className + ".toml"
	conf, err := os.Create(confName)
	if err != nil {
		return err
	}
	defer conf.Close()
	if err := toml.NewEncoder(conf).Encode(configTxt); err != nil {
		return err
	}
	return nil
}

// exportLabels exports all labels of a given class
func exportLabels(tx *sql.Tx, classID int) ([]string, error) {
	stmt, err := tx.Prepare("SELECT label FROM class_label WHERE project_class_id = ? ORDER BY label ASC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var labels []string
	var label string
	for query.Next() {
		query.Scan(&label)
		labels = append(labels, label)
	}
	return labels, nil
}

// exportFolders exports all folders of a given class
func exportFolders(tx *sql.Tx, classID int) (map[string]string, error) {
	stmt, err := tx.Prepare("SELECT target_path, template_name FROM class_folder WHERE project_class_id = ? ORDER BY target_path")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	folders := make(map[string]string)
	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		if _, ok := folders[target]; !ok {
			folders[target] = template
		}
	}
	return folders, nil
}

// exportFiles exports all files of a given class
func exportFiles(tx *sql.Tx, classID int) (map[string]string, error) {
	stmt, err := tx.Prepare("SELECT target_path, template_name FROM class_file WHERE project_class_id = ? ORDER BY target_path")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	files := make(map[string]string)
	var target, template string
	for query.Next() {
		query.Scan(&target, &template)
		if _, ok := files[target]; !ok {
			files[target] = template
		}
	}
	return files, nil
}

// exportScripts exports all scripts of a given class
func exportScripts(tx *sql.Tx, classID int) (map[string]bool, error) {
	stmt, err := tx.Prepare("SELECT script_name, run_as_sudo FROM class_script WHERE project_class_id = ? ORDER BY script_name ASC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	query, err := stmt.Query(classID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	scripts := make(map[string]bool)
	var scriptName string
	var runAsSudo bool
	for query.Next() {
		query.Scan(&scriptName, &runAsSudo)
		if _, ok := scripts[scriptName]; !ok {
			scripts[scriptName] = runAsSudo
		}
	}
	return scripts, nil
}
