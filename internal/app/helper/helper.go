package helper

import (
	"database/sql"
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

	if !queryClassID.Next() {
		return -1, fmt.Errorf("could not find class %s in database", className)
	}

	var classID int
	err = queryClassID.Scan(&classID)
	return classID, err
}
