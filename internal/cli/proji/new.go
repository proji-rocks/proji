package proji

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/plugins"
	"github.com/nikoksr/proji/pkg/templates"
)

// projectNewCommand returns a new instance of the new command.
func projectNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "Create a new project",
		Aliases: []string{"do", "create"},
		Args:    cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			packageLabel := args[0]
			path := args[1]

			return newProject(cmd.Context(), packageLabel, path)
		},
	}

	return cmd
}

func localPathToAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get current working directory")
	}

	return filepath.Join(cwd, path), nil
}

var missingTemplateKeyFn templates.MissingKeyFn = func(key string) (value string, err error) {
	_, err = fmt.Printf("   > %s: ", cases.Title(language.Und, cases.NoLower).String(key))
	if err != nil {
		return "", errors.Wrapf(err, "prompt input for template key %q", key)
	}

	reader := bufio.NewReader(os.Stdin)
	value, err = reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrapf(err, "read input for template key %q", key)
	}

	// Trim newline; note: strings.TrimSuffix checks if the string ends with the suffix before trimming
	value = strings.TrimSuffix(value, "\n")

	return value, nil
}

func createEntry(ctx context.Context, entry *domain.DirEntry, templatesDir string, tmpl *templates.TemplateEngine) error {
	logger := simplog.FromContext(ctx)

	// Check if template path is a template string
	entryPath := entry.Path
	if parsedPath, err := tmpl.ParseToString(ctx, entryPath); err != nil {
		logger.Debugf("template path %q is not a template string", entryPath)
	} else {
		entryPath = parsedPath
	}

	// If we have a file, get its directory and create it. This allows for implicit directory creation and may
	// simplify the directory tree structure in a packages config vastly.
	dirPath := entryPath
	filePath := ""
	if !entry.IsDir {
		dirPath = filepath.Dir(entryPath)
		filePath = entryPath
	}

	if dirPath != "." {
		logger.Debugf("creating directory %q", dirPath)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return errors.Wrapf(err, "create directory %q", dirPath)
		}
	}

	// Skip if we don't have a file path
	if filePath == "" {
		return nil
	}

	logger.Debugf("creating file %q", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		if os.IsExist(err) {
			// If the file already exists, we can ignore the error.
			logger.Warnf("file %q already exists; skipping", filePath)
		} else {
			return errors.Wrapf(err, "create file %q", filePath)
		}
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			logger.Errorf("error closing file %q: %v", filePath, ferr)
		}
	}()

	// Skip if we have no template
	if entry.Template == nil {
		return nil
	}

	tmplPath := entry.Template.Path
	if tmplPath == "" {
		logger.Warnf("template %q has no path; skipping", entry.Template.ID)
		return nil
	}

	// Check if template path is absolute
	if !filepath.IsAbs(tmplPath) {
		tmplPath = filepath.Join(templatesDir, tmplPath)
	}

	// Parse template
	logger.Debugf("generating file %q from template %q", entryPath, entry.Template.ID)
	logger.Debugf("parsing template from file %q", tmplPath)
	if err = tmpl.ParseFile(ctx, file, tmplPath); err != nil {
		return errors.Wrapf(err, "parse template from file %q", tmplPath)
	}

	return nil
}

func runPlugin(ctx context.Context, plugin *domain.Plugin, pluginsDir string) error {
	logger := simplog.FromContext(ctx)

	path := plugin.Path
	if path == "" {
		logger.Warnf("plugin %q has no path; skipping", plugin.ID)
		return nil
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(pluginsDir, path)
	}

	logger.Infof("Running plugin %q", filepath.Base(path))

	return plugins.Run(ctx, path)
}

func buildProject(ctx context.Context, project *domain.ProjectAdd) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	session := cli.SessionFromContext(ctx)

	// Get mandatory paths for finding plugins and templates in the filesystem
	config := session.Config
	pluginsDir := config.PluginsDir()
	templatesDir := config.TemplatesDir()

	// Get package manager from session
	pama := session.PackageManager
	if pama == nil {
		return errors.New("no package manager found")
	}

	// Try to load package by label
	logger.Debugf("loading package %s", project.Package)
	_package, err := pama.GetByLabel(ctx, project.Package)
	if err != nil {
		return errors.Wrapf(err, "get package %q", project.Package)
	}

	// Create project from package at path
	logger.Debugf("creating project from package %q at path %q", _package.Label, project.Path)

	// Create base directory
	logger.Infof("Creating base directory %q", project.Path)
	if err = os.Mkdir(project.Path, 0o755); err != nil {
		if os.IsExist(err) {
			return errors.Newf("path %q already exists", project.Path)
		}

		return errors.Wrapf(err, "create project at path %q", project.Path)
	}

	// Get current working directory
	logger.Debugf("getting current working directory")
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get current working directory")
	}
	logger.Debugf("current working directory: %s", cwd)

	// Change to base path
	logger.Debugf("changing to project base path %q", project.Path)
	if err = os.Chdir(project.Path); err != nil {
		return errors.Wrapf(err, "change to base path %q", project.Path)
	}

	// Make sure to change back to original directory
	defer func() {
		logger.Debugf("changing back to original directory %q", cwd)
		if ferr := os.Chdir(cwd); ferr != nil {
			err = errors.CombineErrors(err, ferr)
		}
	}()

	// Pre-run plugins
	if _package.Plugins != nil {
		for _, plugin := range _package.Plugins.Pre {
			if err = runPlugin(ctx, plugin, pluginsDir); err != nil {
				return errors.Wrapf(err, "run pre-run plugin %q", plugin.ID)
			}
		}
	}

	// Create project in filesystem; meaning file structure and templates
	if _package.DirTree != nil {
		// Create template engine using default tags.
		tmpl := templates.NewEngine("", "")
		tmpl.MissingKeyFn = missingTemplateKeyFn

		logger.Infof("Creating project structure")
		for _, entry := range _package.DirTree.Entries {
			if err = createEntry(ctx, entry, templatesDir, tmpl); err != nil {
				return errors.Wrapf(err, "create directory tree entry %q", entry.Path)
			}
		}
	}

	// Post-run plugins
	if _package.Plugins != nil {
		for _, plugin := range _package.Plugins.Post {
			if err = runPlugin(ctx, plugin, pluginsDir); err != nil {
				return errors.Wrapf(err, "run post-run plugin %q", plugin.ID)
			}
		}
	}

	return nil
}

func newProject(ctx context.Context, packageLabel, name string) error {
	logger := simplog.FromContext(ctx)

	// Get project manager from session
	logger.Debug("getting project manager from cli session")
	prma := cli.SessionFromContext(ctx).ProjectManager
	if prma == nil {
		return errors.New("no project manager found")
	}

	// Get absolute path to project
	name = strings.TrimSpace(name)
	path, err := localPathToAbsPath(name)
	if err != nil {
		return errors.Wrapf(err, "get absolute path to project %q", name)
	}

	// Create project from package at path
	project := domain.NewProject(packageLabel, path, name)

	err = buildProject(ctx, project)
	if err != nil {
		return errors.Wrapf(err, "build project %q at %q from %q", project.Name, project.Path, project.Package)
	}

	// Store project through project manager
	logger.Debug("storing project %q in project manager", project.Name)

	err = prma.Store(ctx, project)
	if err != nil {
		return errors.Wrapf(err, "store project %q in project manager", project.Name)
	}

	logger.Infof("Successfully created project %q", project.Path)

	return nil
}
