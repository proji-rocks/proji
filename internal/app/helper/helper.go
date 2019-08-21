package helper

import "strings"

// ProjectHeader returns an individual graphical header for a project
func ProjectHeader(projectName string) string {
	separatorLine := strings.Repeat("#", 50) + "\n"
	projectLine := "# " + projectName + "\n"
	return (separatorLine + "#\n" + projectLine + "#\n" + separatorLine)
}
