package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"

	"github.com/nikoksr/proji/internal/message"
	"github.com/spf13/cobra"
)

type packageAddCommand struct {
	cmd *cobra.Command
}

func newPackageAddCommand() *packageAddCommand {
	var cmd = &cobra.Command{
		Use:                   "add NAME [NAME...]",
		Short:                 "Add one or more packages",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Deprecated:            "command 'package add' will be deprecated in the next release",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				err := addPackage(name)
				if err != nil {
					message.Warningf("adding package %s failed, %s", name, err.Error())
				} else {
					message.Infof("package %s was successfully added", name)
				}
			}
			return nil
		},
	}

	return &packageAddCommand{cmd: cmd}
}

func addPackage(name string) error {
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

	pkg := domain.NewPackage(name, label)
	pkg.Templates = templates
	pkg.Plugins = plugins

	return session.packageService.StorePackage(pkg)
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

func getTemplates(reader *bufio.Reader) ([]*domain.Template, error) {
	fmt.Println("> Templates (IsFile Destination [Template])")
	templates := make([]*domain.Template, 0)
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
			message.Warningf("minimum of 2 arguments needed")
			continue InputLoop
		case lenInput > 3:
			message.Warningf("more than 3 arguments given")
			continue InputLoop
		}

		isFile, err := strconv.ParseBool(splittedInput[0])
		if err != nil {
			message.Warningf("value given for 'IsFile' field is not a boolean (true|false)")
			continue InputLoop
		}
		destination := splittedInput[1]

		// Check if dest exists
		if _, ok := destinations[destination]; ok {
			message.Warningf("destination path %s was already defined", destination)
			continue InputLoop
		}

		path := ""
		if len(splittedInput) == 3 {
			path = splittedInput[2]
		}

		// Add folder(s) to map
		destinations[destination] = true
		templates = append(templates, &domain.Template{
			IsFile:      isFile,
			Path:        path,
			Destination: destination,
		})
	}
	fmt.Println()
	return templates, nil
}

func getPlugins(reader *bufio.Reader) ([]*domain.Plugin, error) {
	fmt.Println("> Plugins (Name Path Execution Number)")
	plugins := make([]*domain.Plugin, 0)
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
			message.Warningf("minimum of 3 arguments needed")
			continue
		}

		pluginPath := splittedInput[0]
		if _, ok := pluginPaths[pluginPath]; ok {
			fmt.Printf("> Warning: Script %s is already in execution list\n", pluginPath)
			continue
		}
		pluginPaths[pluginPath] = true

		execNumber, err := strconv.Atoi(splittedInput[1])
		if err != nil {
			message.Warningf("value given for 'ExecNumber' field is not a integer")
			continue
		}
		if execNumber == 0 {
			message.Warningf("execution number may not be equal to zero")
			continue
		}
		if _, ok := execNumbers[execNumber]; ok {
			message.Warningf("execution number %d was already given", execNumber)
			continue
		}
		execNumbers[execNumber] = true

		// Get optional description
		var pluginDescription string
		if len(splittedInput) == 3 {
			pluginDescription = splittedInput[2]
		}
		plugins = append(plugins, &domain.Plugin{
			Path:        pluginPath,
			ExecNumber:  execNumber,
			Description: pluginDescription,
		})
	}
	fmt.Println()
	return plugins, nil
}
