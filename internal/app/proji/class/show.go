package class

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/helper"
)

// Show shows detailed information abour a given class
func (c *Class) Show() error {
	if err := c.Load(); err != nil {
		return err
	}

	fmt.Println(helper.ProjectHeader(c.Name))
	c.showLabels()
	c.showFolders()
	c.showFiles()
	c.showScripts()
	return nil
}

// showLabels shows all labels of a given class
func (c *Class) showLabels() {
	fmt.Println("Labels:")

	for _, label := range c.Labels {
		fmt.Println(" " + label)
	}
	fmt.Println()
}

// showFolders shows all folders of a given class
func (c *Class) showFolders() {
	fmt.Println("Folders:")

	for target, template := range c.Folders {
		fmt.Println(" " + target + " : " + template)
	}
	fmt.Println()
}

// showFiles shows all files of a given class
func (c *Class) showFiles() {
	fmt.Println("Files:")

	for target, template := range c.Files {
		fmt.Println(" " + target + " : " + template)
	}
	fmt.Println()
}

// showScripts shows all scripts of a given class
func (c *Class) showScripts() {
	fmt.Println("Scripts:")

	for script, runAsSudo := range c.Scripts {
		sudo := ""
		if runAsSudo {
			sudo = "sudo "
		}
		fmt.Println(" " + sudo + script)
	}
	fmt.Println()
}
