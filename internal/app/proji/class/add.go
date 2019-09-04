package class

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Add adds a new class interactively through the cli.
func Add(name string) (*Class, error) {
	c := New(name)
	reader := bufio.NewReader(os.Stdin)

	if err := c.addLabels(reader); err != nil {
		return nil, err
	}
	if err := c.addFolders(reader); err != nil {
		return nil, err
	}
	if err := c.addFiles(reader); err != nil {
		return nil, err
	}
	if err := c.addScripts(reader); err != nil {
		return nil, err
	}

	return c, c.Save()
}

// addLabels adds labels to the class.
func (c *Class) addLabels(reader *bufio.Reader) error {
	fmt.Print("Labels: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	c.Labels = strings.Fields(text)
	if len(c.Labels) < 1 {
		return fmt.Errorf("you have to specify atleast one label")
	}

	fmt.Println()
	return nil
}

// addFolders adds folders to the class.
func (c *Class) addFolders(reader *bufio.Reader) error {
	fmt.Println("Folders: ")

	for {
		// Read in folders
		// Syntax: Target [Source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		folderPair := strings.Fields(input)
		numFolders := len(folderPair)

		// End if no input given
		if numFolders < 1 {
			break
		}
		if numFolders > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		target := folderPair[0]

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		if src, ok := c.Folders[target]; ok {
			fmt.Printf("Warning: Target folder %s is already associated to source folder %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFolders > 1 {
			src = folderPair[1]
		}

		// Add folder(s) to map
		c.Folders[target] = src
	}
	fmt.Println()
	return nil
}

// addFiles adds files to the class.
func (c *Class) addFiles(reader *bufio.Reader) error {
	fmt.Println("Files: ")

	for {
		// Read in files
		// Syntax: target [source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		filePair := strings.Fields(input)
		numFiles := len(filePair)

		// End if no input given
		if numFiles < 1 {
			break
		}
		if numFiles > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		target := filePair[0]

		if src, ok := c.Files[target]; ok {
			fmt.Printf("Warning: Target file %s is already associated to source file %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFiles > 1 {
			src = filePair[1]
		}

		// Add file(s) to map
		c.Files[target] = src
	}
	fmt.Println()
	return nil
}

// addScripts adds scripts to the class.
func (c *Class) addScripts(reader *bufio.Reader) error {
	fmt.Println("Scripts: ")

	for {
		// Read in scripts
		// Syntax: script [sudo]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		scriptData := strings.Fields(input)
		lenData := len(scriptData)

		// End if no input given
		if lenData < 1 {
			break
		}
		if lenData > 2 {
			fmt.Println("Warning: More than two files were given.")
			continue
		}

		// Set sudo to true if given
		var sudo bool
		script := scriptData[0]

		if lenData == 2 {
			if scriptData[0] != "sudo" {
				fmt.Printf("Warning: %s invalid. Has to be 'sudo' or ''(empty).", scriptData[0])
			}
			sudo = true
			script = scriptData[1]
		}

		// Check if target exists
		// A target should only exist once
		// A source can be used multiple times
		if _, ok := c.Scripts[script]; ok {
			fmt.Printf("Warning: Script %s is already in execution list.\n", script)
			continue
		}

		// Add folder(s) to map
		c.Scripts[script] = sudo
	}
	fmt.Println()
	return nil
}
