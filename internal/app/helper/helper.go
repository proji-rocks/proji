package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ProjectHeader returns an individual graphical header for a project
func ProjectHeader(projectName string) string {
	separatorLine := strings.Repeat("#", 50) + "\n"
	projectLine := "# " + projectName + "\n"
	return (separatorLine + "#\n" + projectLine + "#\n" + separatorLine)
}

// ParseArgs parses the cli arguments to the needed data - the extension and the project names.
// AreArgsValid() should be run before this function.
func ParseArgs() (string, []string, error) {
	args := os.Args[1:]

	if len(args) < 2 {
		return "", []string{}, errors.New("insufficient number of cli arguments")
	}

	return args[0], args[1:], nil
}

// GetConfigDir returns the default config directory.
func GetConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home + "/.config/proji/"
}

// QueryClassID queries the id of namely specified class
func QueryClassID(tx *sql.Tx, className string) (int, error) {
	stmt, err := tx.Prepare("SELECT class_id FROM class WHERE name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	queryClassID, err := stmt.Query(className)
	if err != nil {
		return -1, err
	}
	defer queryClassID.Close()

	var classID int
	if !queryClassID.Next() {
		return -1, fmt.Errorf("could not find class %s in database", className)
	}
	err = queryClassID.Scan(&classID)
	if err != nil {
		return -1, err
	}
	return classID, nil
}
