package helper

import (
	"errors"
	"os"
	"strings"
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
		return "", []string{}, errors.New("Insufficient number of cli arguments.\nSyntax: create_project <extension> <project> [projects]")
	}

	return args[0], args[1:], nil
}
