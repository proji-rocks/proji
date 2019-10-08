package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProjectHeader returns an individual graphical header for a project
func ProjectHeader(title string) string {
	numChars := len(title) + 4
	separatorLine := "+" + strings.Repeat("-", numChars-2) + "+\n"
	projectLine := "| " + title + " |\n"
	return (separatorLine + projectLine + separatorLine)
}

// DoesPathExist checks if a given path exists in the filesystem.
func DoesPathExist(path string) bool {
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
