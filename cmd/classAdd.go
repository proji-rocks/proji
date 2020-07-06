//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nikoksr/proji/storage/models"

	"github.com/spf13/cobra"
)

var classAddCmd = &cobra.Command{
	Use:        "add NAME [NAME...]",
	Short:      "Add one or more classes",
	Deprecated: "command 'class add' will be deprecated in the next release",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing class name")
		}

		for _, name := range args {
			err := addClass(name)
			if err != nil {
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

func addClass(name string) error {
	reader := bufio.NewReader(os.Stdin)

	label, err := getLabel(reader)
	if err != nil {
		return err
	}
	templates, err := getTemplates(reader)
	if err != nil {
		return err
	}
	plugins, err := getPlugins(reader)
	if err != nil {
		return err
	}

	class := models.NewClass(name, label, false)
	class.Templates = templates
	class.Plugins = plugins
	return projiEnv.StorageService.SaveClass(class)
}

func getLabel(reader *bufio.Reader) (string, error) {
	fmt.Print("> Label: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	labels := strings.Fields(text)
	if len(labels) > 1 {
		return "", fmt.Errorf("only one label is needed")
	}
	fmt.Println()
	return labels[0], nil
}

func getTemplates(reader *bufio.Reader) ([]*models.Template, error) {
	fmt.Println("> Templates (IsFile Destination [Template]): ")
	templates := make([]*models.Template, 0)
	destinations := make(map[string]bool)

InputLoop:
	for {
		// Read in folders
		// Syntax: IsFile Destination [template]
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		splittedInput := strings.Fields(input)
		lenInput := len(splittedInput)

		// End if no input given
		switch {
		case lenInput < 1:
			break InputLoop
		case lenInput < 2:
			fmt.Println("> Warning: At least 2 arguments needed.")
			continue InputLoop
		case lenInput > 3:
			fmt.Println("> Warning: More than three arguments were given.")
			continue InputLoop
		}

		isFile, err := strconv.ParseBool(splittedInput[0])
		if err != nil {
			fmt.Println("> Warning: Value given for 'IsFile' field is not a boolean (true|false).")
			continue InputLoop
		}
		destination := splittedInput[1]

		// Check if dest exists
		if _, ok := destinations[destination]; ok {
			fmt.Printf("> Warning: Destination path '%s' was already defined.\n", destination)
			continue InputLoop
		}

		path := ""
		if len(splittedInput) == 3 {
			path = splittedInput[2]
		}

		// Add folder(s) to map
		destinations[destination] = true
		templates = append(templates, &models.Template{
			IsFile:      isFile,
			Path:        path,
			Destination: destination,
		})
	}
	fmt.Println()
	return templates, nil
}

func getPlugins(reader *bufio.Reader) ([]*models.Plugin, error) {
	fmt.Println("> Plugins (Name Path Execution Number): ")
	plugins := make([]*models.Plugin, 0)
	pluginPaths := make(map[string]bool)
	execNumbers := make(map[int]bool)

	for {
		// Read in plugins
		// Syntax: name path execNumber
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		splittedInput := strings.Fields(input)
		lenInput := len(splittedInput)

		// End if no input given
		if lenInput < 1 {
			break
		} else if lenInput < 3 {
			fmt.Println("> Warning: 3 arguments needed.")
			continue
		}

		pluginName := splittedInput[0]
		pluginPath := splittedInput[1]
		if _, ok := pluginPaths[pluginPath]; ok {
			fmt.Printf("> Warning: Script %s is already in execution list\n", pluginPath)
			continue
		}
		pluginPaths[pluginPath] = true

		execNumber, err := strconv.Atoi(splittedInput[2])
		if err != nil {
			fmt.Println("> Warning: Value given for 'ExecNumber' field is not a integer.")
			continue
		}
		if execNumber == 0 {
			fmt.Println("> Warning: Execution number may not be equal to zero")
			continue
		}
		if _, ok := execNumbers[execNumber]; ok {
			fmt.Printf("> Warning: Execution number %d was already given\n", execNumber)
			continue
		}
		execNumbers[execNumber] = true
		plugins = append(plugins, &models.Plugin{
			Name:       pluginName,
			Path:       pluginPath,
			ExecNumber: execNumber,
		})
	}
	fmt.Println()
	return plugins, nil
}
