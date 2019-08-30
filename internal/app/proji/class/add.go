package class

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/viper"
)

// AddClassCLI adds a new class interactively through the cli to the database
func AddClassCLI(className string) error {
	className = strings.ToLower(className)
	reader := bufio.NewReader(os.Stdin)

	labels, err := addLabels(reader)
	if err != nil {
		return err
	}
	folders, err := addFolders(reader)
	if err != nil {
		return err
	}
	files, err := addFiles(reader)
	if err != nil {
		return err
	}
	scripts, err := addScripts(reader)
	if err != nil {
		return err
	}

	err = AddClassToDB(className, labels, folders, files, scripts)
	if err != nil {
		return err
	}

	fmt.Printf("> Added class %s successfully.\n", className)
	return nil
}

// addLabels adds the labels related to the new class
func addLabels(reader *bufio.Reader) ([]string, error) {
	fmt.Print("Labels: ")
	text, err := reader.ReadString('\n')

	if err != nil {
		return []string{}, err
	}

	labels := strings.Fields(text)

	if len(labels) < 1 {
		return labels, fmt.Errorf("you have to specify atleast one label")
	}

	fmt.Println()
	return labels, nil
}

// addFiles adds the files related to the new class
func addFiles(reader *bufio.Reader) (map[string]string, error) {
	fmt.Println("Files: ")
	allFiles := make(map[string]string)

	for {
		// Read in files
		// Syntax: target [source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		files := strings.Fields(input)

		if err != nil {
			return map[string]string{}, err
		}

		numFiles := len(files)

		// End if no input given
		if numFiles < 1 {
			break
		}
		if numFiles > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		target := files[0]

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		if _, ok := allFiles[target]; ok {
			fmt.Printf("Warning: Target file %s is already associated to a source file.\n", target)
			continue
		}

		// Add source if given
		source := ""

		if numFiles > 1 {
			source = files[1]
		}

		// Add file(s) to map
		allFiles[target] = source
	}
	fmt.Println()
	return allFiles, nil
}

// addFolders adds the folders related to the new class
func addFolders(reader *bufio.Reader) (map[string]string, error) {
	fmt.Println("Folders: ")
	allFolders := make(map[string]string)

	for {
		// Read in folders
		// Syntax: Target [Source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		folders := strings.Fields(input)

		if err != nil {
			return map[string]string{}, err
		}

		numFolders := len(folders)

		// End if no input given
		if numFolders < 1 {
			break
		}
		if numFolders > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		target := folders[0]

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		if _, ok := allFolders[target]; ok {
			fmt.Printf("Warning: Target folder %s is already associated to a source folder.\n", target)
			continue
		}

		// Add source if given
		source := ""

		if numFolders > 1 {
			source = folders[1]
		}

		// Add folder(s) to map
		allFolders[target] = source
	}
	fmt.Println()
	return allFolders, nil
}

// addScripts adds the scripts related to the new class
func addScripts(reader *bufio.Reader) (map[string]bool, error) {
	fmt.Println("Scripts: ")
	allScripts := make(map[string]bool)

	for {
		// Read in scripts
		// Syntax: script [sudo]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		scripts := strings.Fields(input)

		if err != nil {
			return map[string]bool{}, err
		}

		numScripts := len(scripts)

		// End if no input given
		if numScripts < 1 {
			break
		}
		if numScripts > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		script := scripts[0]

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		if _, ok := allScripts[script]; ok {
			fmt.Printf("Warning: Script %s is already in excution list.\n", script)
			continue
		}

		// Add source if given
		sudo := false

		if numScripts > 1 && scripts[1] == "sudo" {
			sudo = true
		}

		// Add folder(s) to map
		allScripts[script] = sudo
	}
	fmt.Println()
	return allScripts, nil
}

// AddClassToDB adds a new class and its dependencies to a database
func AddClassToDB(className string, labels []string, folders, files map[string]string, scripts map[string]bool) error {
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

	// Insert new class
	err = insertClass(tx, className)
	if err != nil {
		return err
	}

	// Get id of new class
	classID, err := helper.QueryClassID(tx, className)
	if err != nil {
		return err
	}

	// Insert class labels
	err = insertLabels(tx, classID, labels)
	if err != nil {
		return err
	}

	// Insert class folders
	err = insertFolders(tx, classID, folders)
	if err != nil {
		return err
	}

	// Insert class files
	err = insertFiles(tx, classID, files)
	if err != nil {
		return err
	}

	// Insert class scripts
	err = insertScripts(tx, classID, scripts)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// insertClass inserts a new class name into the database
func insertClass(tx *sql.Tx, className string) error {
	stmt, err := tx.Prepare("INSERT INTO class(name) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(className)
	if err != nil {
		return err
	}
	return nil
}

// insertLabels inserts new class labels into the database
func insertLabels(tx *sql.Tx, classID int, labels []string) error {
	stmt, err := tx.Prepare("INSERT INTO class_label(class_id, label) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, label := range labels {
		_, err = stmt.Exec(classID, strings.ToLower(label))
		if err != nil {
			return err
		}
	}
	return nil
}

// insertFolders inserts new class folders into the database
func insertFolders(tx *sql.Tx, classID int, folders map[string]string) error {
	stmt, err := tx.Prepare("INSERT INTO class_folder(class_id, target, template) VALUES(?, ?, ?)")
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

// insertFiles inserts new class files into the database
func insertFiles(tx *sql.Tx, classID int, files map[string]string) error {
	stmt, err := tx.Prepare("INSERT INTO class_file(class_id, target, template) VALUES(?, ?, ?)")
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

// insertScripts inserts new class scripts into the database
func insertScripts(tx *sql.Tx, classID int, scripts map[string]bool) error {
	stmt, err := tx.Prepare("INSERT INTO class_script(class_id, script_name, run_as_sudo) VALUES(?, ?, ?)")
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
