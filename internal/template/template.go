package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
	"github.com/valyala/fasttemplate"
)

type engine struct {
	startTag string
	endTag   string
	seenTags map[string]string
}

// CreateTemplatesInProject creates given templates inside of a project directory. It checks if a template has
// a template file assigned to it and if that is the case, it parses the template file and writes it into the project
// directory. When no template file is assigned it just creates an empty file or folder at the destination path
// specified by the template.
func CreateTemplatesInProject(templatesBasePath, projectPath string, templates []*domain.Template) error {
	return createTemplatesInProject(templatesBasePath, projectPath, templates)
}

func createTemplatesInProject(templatesBasePath, projectPath string, templates []*domain.Template) error {
	var err error
	e := engine{
		startTag: "{{%",
		endTag:   "%}}",
		seenTags: map[string]string{
			"project-name": filepath.Base(projectPath),
		},
	}

	for _, template := range templates {
		if isEmptyTemplate(template) {
			err = e.createEmptyTemplate(template)
			if err != nil {
				return errors.Wrap(err, "create empty template")
			}
			continue
		}
		if template.IsFile {
			err = e.parseTemplateFile(
				filepath.Join(templatesBasePath, template.Path),
				filepath.Join(projectPath, template.Destination),
			)
			if err != nil {
				return errors.Wrap(err, "parse template file")
			}
		} else {
			err = e.parseTemplateFolder(
				filepath.Join(templatesBasePath, template.Path),
				filepath.Join(projectPath, template.Destination),
			)
			if err != nil {
				return errors.Wrap(err, "parse template directory")
			}
		}
	}
	return nil
}

func (e *engine) createEmptyTemplate(template *domain.Template) error {
	var err error
	if template.IsFile {
		// Create file
		_, err = os.Create(template.Destination)
	} else {
		// Create folder
		err = os.MkdirAll(template.Destination, os.ModePerm)
	}
	return err
}

func (e *engine) parseTemplateFile(path, destination string) error {
	// Load the template file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse the template
	t, err := fasttemplate.NewTemplate(string(b), e.startTag, e.endTag)
	if err != nil {
		return errors.Wrap(err, "parse template")
	}
	s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		// Check if space holder was already replaced once
		tag = normalizePlaceholder(tag)
		value, exists := e.seenTags[tag]
		if !exists {
			// If not found, read in a replacement for the placeholder
			_, err = fmt.Printf("%s: ", strings.Title(tag))
			if err != nil {
				return -1, errors.Wrap(err, "placeholder prompt")
			}
			_, err = fmt.Scanf("%s", &value)
			if err != nil {
				return -1, errors.Wrap(err, "read placeholder replacement")
			}
			e.seenTags[tag] = value
		}
		return w.Write([]byte(value))
	})

	// Create file from the parsed template in the project folder
	f, err := os.Create(destination)
	if err != nil {
		return errors.Wrap(err, "create file from parsed template")
	}
	defer f.Close()
	_, err = f.Write([]byte(s))
	return err
}

func (e *engine) parseTemplateFolder(path, destination string) error {
	return filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip base directory
		if path == currentPath {
			return nil
		}
		// Extract relative path
		relPath, err := filepath.Rel(path, currentPath)
		if err != nil {
			return err
		}

		// Add file or folder to package
		isFile := true
		if info.IsDir() {
			isFile = false
		}

		// Parse the template file or create the folder
		absoluteDestination := filepath.Join(destination, relPath)
		if isFile {
			err = os.MkdirAll(filepath.Dir(absoluteDestination), os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "create directory of template file")
			}
			err = e.parseTemplateFile(currentPath, absoluteDestination)
		} else {
			err = os.MkdirAll(absoluteDestination, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "create template directory")
			}
		}
		return err
	})
}

// isEmptyTemplate checks if a template object has an actual template assigned to it.
func isEmptyTemplate(template *domain.Template) bool {
	return len(template.Path) == 0
}

// normalizePlaceholder normalizes a placeholder in order to avoid as many duplicate placeholder entries as possible.
// For example, we don't want the map of replaced placeholders to hold an entries for 'project-name', 'projectname' and
// 'Project-Name'. Not only would this result in unnecessarily allocated memory but also in duplicate user input prompts;
// asking three times for the project name would be pretty annoying.
func normalizePlaceholder(placeholder string) string {
	return strings.Trim(strings.ToLower(placeholder), " _-")
}
