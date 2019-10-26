package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
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
			if err := addClass(name, projiEnv.Svc); err != nil {
				fmt.Printf("> Adding class '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("> Class '%s' was successfully added\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classAddCmd)
}

func addClass(name string, svc storage.Service) error {
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

	class := item.NewClass(name, label, false)
	class.Folders = folders
	class.Files = files
	class.Scripts = scripts
	return svc.SaveClass(class)
}

func getLabel(reader *bufio.Reader) (string, error) {
	fmt.Print("> Label: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	labels := strings.Fields(text)
	if len(labels) > 1 {
		return "", fmt.Errorf("Only one label is needed")
	}
	fmt.Println()
	return labels[0], nil
}

func getFolders(reader *bufio.Reader) ([]*item.Folder, error) {
	fmt.Println("> Folders: ")
	folders := make([]*item.Folder, 0)
	destinations := make(map[string]bool)

	for {
		// Read in folders
		// Syntax: Destination [template]
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
			fmt.Println("> Warning: More than two files were given.")
			continue
		}

		dest := folderPair[0]

		// Check if dest exists
		if _, ok := destinations[dest]; ok {
			fmt.Printf("> Warning: Destination folder '%s' has already been defined.\n", dest)
			continue
		}

		// Add template if given
		template := ""
		if numFolders > 1 {
			template = folderPair[1]
		}

		// Add folder(s) to map
		destinations[dest] = true
		folders = append(folders, &item.Folder{Destination: dest, Template: template})
	}
	fmt.Println()
	return folders, nil
}

func getFiles(reader *bufio.Reader) ([]*item.File, error) {
	fmt.Println("> Files: ")
	files := make([]*item.File, 0)
	destinations := make(map[string]bool)

	for {
		// Read in files
		// Syntax: dest [template]
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
			fmt.Println("> Warning: More than two files were given.")
			continue
		}

		// Check if dest exists
		// A dest should only exist once
		// A template can be used multiple times
		dest := filePair[0]

		if _, ok := destinations[dest]; ok {
			fmt.Printf("> Warning: Destination folder '%s' has already been defined.\n", dest)
			continue
		}

		// Add template if given
		template := ""
		if numFiles > 1 {
			template = filePair[1]
		}

		// Add file(s) to map
		destinations[dest] = true
		files = append(files, &item.File{Destination: dest, Template: template})
	}
	fmt.Println()
	return files, nil
}

func getScripts(reader *bufio.Reader) ([]*item.Script, error) {
	fmt.Println("> Scripts: ")
	scripts := make([]*item.Script, 0)
	scriptNames := make(map[string]bool)

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
			fmt.Println("> Warning: More than two files were given.")
			continue
		}

		// Set sudo to true if given
		var runAsSudo bool
		scriptName := scriptData[0]

		if lenData == 2 {
			if scriptData[0] != "sudo" {
				fmt.Printf("> Warning: %s invalid. Has to be 'sudo' or ''(empty).", scriptData[0])
			}
			runAsSudo = true
			scriptName = scriptData[1]
		}

		// Check if script was already added to execution list
		if _, ok := scriptNames[scriptName]; ok {
			fmt.Printf("> Warning: Script %s is already in execution list\n", scriptName)
			continue
		}

		scriptNames[scriptName] = true
		scripts = append(scripts, &item.Script{Name: scriptName, RunAsSudo: runAsSudo})
	}
	fmt.Println()
	return scripts, nil
}
