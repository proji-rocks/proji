package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/proji/storage"

	"github.com/nikoksr/proji/pkg/helper"

	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:   "add NAME [NAME...]",
	Short: "Add one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing class name")
		}

		for _, name := range args {
			if err := AddClass(name); err != nil {
				fmt.Printf("Adding class '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("Class '%s' was successfully added.\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}

// AddClass adds a new class interactively through the cli.
func AddClass(name string) error {
	reader := bufio.NewReader(os.Stdin)

	label, err := getLabel(reader)
	if err != nil {
		return err
	}
	folders, err := getFolders(reader)
	if err != nil {
		return err
	}
	files, err := getFiles(reader)
	if err != nil {
		return err
	}
	scripts, err := getScripts(reader)
	if err != nil {
		return err
	}

	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	class, err := storage.NewClass(name, label)
	if err != nil {
		return err
	}
	class.Folders = folders
	class.Files = files
	class.Scripts = scripts
	return s.SaveClass(class)
}

func getLabel(reader *bufio.Reader) (string, error) {
	fmt.Print("Label: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	labels := strings.Fields(text)
	if len(labels) > 1 {
		return "", fmt.Errorf("Only one unique label is allowed")
	}
	fmt.Println()
	return labels[0], nil
}

func getFolders(reader *bufio.Reader) (map[string]string, error) {
	fmt.Println("Folders: ")
	folders := make(map[string]string)

	for {
		// Read in folders
		// Syntax: Target [Source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
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
		if src, ok := folders[target]; ok {
			fmt.Printf("Warning: Target folder %s is already associated to source folder %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFolders > 1 {
			src = folderPair[1]
		}

		// Add folder(s) to map
		folders[target] = src
	}
	fmt.Println()
	return folders, nil
}

func getFiles(reader *bufio.Reader) (map[string]string, error) {
	fmt.Println("Files: ")
	files := make(map[string]string)

	for {
		// Read in files
		// Syntax: target [source]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
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

		if src, ok := files[target]; ok {
			fmt.Printf("Warning: Target file %s is already associated to source file %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFiles > 1 {
			src = filePair[1]
		}

		// Add file(s) to map
		files[target] = src
	}
	fmt.Println()
	return files, nil
}

func getScripts(reader *bufio.Reader) (map[string]bool, error) {
	fmt.Println("Scripts: ")
	scripts := make(map[string]bool)

	for {
		// Read in scripts
		// Syntax: script [sudo]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
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
		if _, ok := scripts[script]; ok {
			fmt.Printf("Warning: Script %s is already in execution list.\n", script)
			continue
		}

		// Add folder(s) to map
		scripts[script] = sudo
	}
	fmt.Println()
	return scripts, nil
}
