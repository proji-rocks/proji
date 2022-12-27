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
	err = os.Mkdir(project.Path, 0o755)
	if err != nil {
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
	err = os.Chdir(project.Path)
	if err != nil {
		return errors.Wrapf(err, "change to base path %q", project.Path)
	}

	// Make sure to change back to original directory
	defer func() {
		logger.Debugf("changing back to original directory %q", cwd)
		ferr := os.Chdir(cwd)
		if ferr != nil {
			err = errors.CombineErrors(err, ferr)
		}
	}()

	// Pre-run plugins
	if _package.Plugins != nil {
		for _, plugin := range _package.Plugins.Pre {
			path := plugin.Path
			if path == "" {
				logger.Warnf("plugin %q has no path; skipping", plugin.ID)
				continue
			}

			if !filepath.IsAbs(path) {
				path = filepath.Join(pluginsDir, path)
			}

			logger.Infof("Running plugin %q", filepath.Base(path))
			err = plugins.Run(ctx, path)
			if err != nil {
				return errors.Wrapf(err, "run plugin %q at %q", plugin.ID, path)
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
			// Check if template path is a template string
			entryPath := entry.Path
			if parsedPath, err := tmpl.ParseString(ctx, entryPath); err != nil {
				logger.Debugf("template path %q is not a template string", entryPath)
			} else {
				entryPath = parsedPath
			}

			if entry.Template != nil {
				logger.Debugf("generating file %q from template %q", entryPath, entry.Template.ID)

				tmplPath := entry.Template.Path
				if tmplPath == "" {
					logger.Warnf("template %q has no path; skipping", entry.Template.ID)
					continue
				}

				// Check if template path is absolute
				if !filepath.IsAbs(tmplPath) {
					tmplPath = filepath.Join(templatesDir, tmplPath)
				}

				// Read template from filesystem
				logger.Debugf("loading template file %q", tmplPath)
				templateData, err := os.ReadFile(tmplPath)
				if err != nil {
					return errors.Wrapf(err, "load template file %q", tmplPath)
				}

				// Create empty destination file
				destinationFile, err := os.Create(entryPath)
				if err != nil {
					return errors.Wrapf(err, "create template destination file %q", entryPath)
				}

				err = tmpl.Parse(ctx, destinationFile, templateData)
				if err != nil {
					return errors.Wrapf(err, "parse template %q", tmplPath)
				}

				continue
			}

			if entry.IsDir {
				logger.Debugf("creating directory %q", entryPath)
				err = os.MkdirAll(entryPath, 0o755)
				if err != nil {
					return errors.Wrapf(err, "create directory %q", entryPath)
				}
			} else {
				logger.Debugf("creating file %q", entryPath)
				_, err = os.Create(entryPath)
				if err != nil {
					return errors.Wrapf(err, "create file %q", entryPath)
				}
			}
		}
	}

	// Post-run plugins
	if _package.Plugins != nil {
		for _, plugin := range _package.Plugins.Post {
			path := plugin.Path
			if path == "" {
				logger.Warnf("plugin %q has no path; skipping", plugin.ID)
				continue
			}

			if !filepath.IsAbs(path) {
				path = filepath.Join(pluginsDir, path)
			}

			logger.Infof("Running plugin %q", filepath.Base(path))
			err = plugins.Run(ctx, path)
			if err != nil {
				return errors.Wrapf(err, "run plugin %q at %q", plugin.ID, path)
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
