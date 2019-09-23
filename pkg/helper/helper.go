package helper

import (
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
