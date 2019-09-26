package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// ProjectHeader returns an individual graphical header for a project
func ProjectHeader(title string) string {
	numChars := len(title) + 4
	separatorLine := "+" + strings.Repeat("-", numChars-2) + "+\n"
	projectLine := "| " + title + " |\n"
	return (separatorLine + projectLine + separatorLine)
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

// GetSqlitePath returns the default location for the sqlite3 db.
func GetSqlitePath() (string, error) {
	dbPath, ok := viper.Get("sqlite3.path").(string)
	if !ok {
		return "", fmt.Errorf("Could not read database name from config file")
	}

	return GetConfigDir() + dbPath, nil
}

// DoesFileExist checks if a given file exists.
func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// StrToUInt converts a string into a uint.
func StrToUInt(num string) (uint, error) {
	// Parse the input
	id64, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

// WantTo waits for a valid user input to confirm if he wants to do whatever was asked for.
func WantTo(question string) bool {
	// Ask to replace project
	var input string
	for {
		fmt.Print(question + " [y/N] ")
		n, err := fmt.Scan(&input)
		if n == 1 && err == nil {
			input = strings.ToLower(input)
			if input == "n" || input == "\n" {
				return false
			}
			if input == "y" {
				return true
			}
		}
	}
}
