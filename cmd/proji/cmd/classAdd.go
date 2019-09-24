package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:   "add CLASS [CLASS...]",
	Short: "add new classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}

		for _, name := range args {
			if err := AddClass(name); err != nil {
				fmt.Printf("Failed adding class %s: %v\n", name, err)
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}

// AddClass adds a new class interactively through the cli.
func AddClass(name string) error {
	// Create class and fill it with data from the cli
	c := storage.NewClass(name)

	reader := bufio.NewReader(os.Stdin)

	if err := getLabels(reader, c); err != nil {
		return err
	}
	if err := getFolders(reader, c); err != nil {
		return err
	}
	if err := getFiles(reader, c); err != nil {
		return err
	}
	if err := getScripts(reader, c); err != nil {
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
	return s.SaveClass(c)
}

func getLabels(reader *bufio.Reader, class *storage.Class) error {
	fmt.Print("Labels: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	class.Labels = strings.Fields(text)
	if len(class.Labels) < 1 {
		return fmt.Errorf("you have to specify atleast one label")
	}

	fmt.Println()
	return nil
}

func getFolders(reader *bufio.Reader, class *storage.Class) error {
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
		if src, ok := class.Folders[target]; ok {
			fmt.Printf("Warning: Target folder %s is already associated to source folder %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFolders > 1 {
			src = folderPair[1]
		}

		// Add folder(s) to map
		class.Folders[target] = src
	}
	fmt.Println()
	return nil
}

func getFiles(reader *bufio.Reader, class *storage.Class) error {
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

		if src, ok := class.Files[target]; ok {
			fmt.Printf("Warning: Target file %s is already associated to source file %s.\n", target, src)
			continue
		}

		// Add source if given
		src := ""
		if numFiles > 1 {
			src = filePair[1]
		}

		// Add file(s) to map
		class.Files[target] = src
	}
	fmt.Println()
	return nil
}

func getScripts(reader *bufio.Reader, class *storage.Class) error {
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
		if _, ok := class.Scripts[script]; ok {
			fmt.Printf("Warning: Script %s is already in execution list.\n", script)
			continue
		}

		// Add folder(s) to map
		class.Scripts[script] = sudo
	}
	fmt.Println()
	return nil
}
