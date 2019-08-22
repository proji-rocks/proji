package main

import (
	"fmt"
	"os"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/nikoksr/proji/internal/app/proji"
)

func main() {
	// Parse the arguments extension and project names
	ext, projects, err := helper.ParseArgs()

	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Load values from a config file
	homeDir := os.Getenv("HOME")
	configDir := homeDir + "/.config/create_project/"
	databaseName := "cp.sqlite3"

	// Get current working directory
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create setup
	newSetup := proji.Setup{Owd: cwd, ConfigDir: configDir, DatabaseName: databaseName, Extension: ext}
	err = newSetup.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newSetup.Stop()

	// Projects loop
	for _, project := range projects {
		fmt.Println(helper.ProjectHeader(project))
		newProject := proji.Project{Name: project, Data: &newSetup}
		err = newProject.Create()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
