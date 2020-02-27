package item

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/repo"
	"github.com/nikoksr/proji/pkg/proji/repo/github"
	"github.com/nikoksr/proji/pkg/proji/repo/gitlab"

	"github.com/BurntSushi/toml"
)

// Class struct represents a proji class
type Class struct {
	ID        uint      // Class ID in storage. Not exported/imported.
	IsDefault bool      // Is this a default class? Not exported/imported.
	Name      string    `toml:"name"`   // Class name
	Label     string    `toml:"label"`  // Class label
	Folders   []*Folder `toml:"folder"` // Class folders
	Files     []*File   `toml:"file"`   // Class files
	Scripts   []*Script `toml:"script"` // Class scripts
}

// NewClass returns a new class
func NewClass(name, label string, isDefault bool) *Class {
	return &Class{
		ID:        0,
		IsDefault: isDefault,
		Name:      name,
		Label:     label,
		Folders:   make([]*Folder, 0),
		Files:     make([]*File, 0),
		Scripts:   make([]*Script, 0),
	}
}

// ImportFromConfig imports class data from a given config file.
func (c *Class) ImportFromConfig(configName string) error {
	// Validate that it's a toml file
	if !strings.HasSuffix(configName, ".toml") {
		return fmt.Errorf("import file has to be of type 'toml'")
	}

	// Validate config is not empty
	conf, err := os.Stat(configName)
	if err != nil {
		return err
	}
	if conf.Size() == 0 {
		return fmt.Errorf("import file is empty")
	}

	// Decode the file
	_, err = toml.DecodeFile(configName, &c)
	if err != nil {
		return err
	}

	if len(c.Name) < 1 {
		return fmt.Errorf("name cannot be an empty string")
	}
	if len(c.Label) < 1 {
		return fmt.Errorf("label cannot be an empty string")
	}

	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Config might be empty")
	}
	return nil
}

// ImportFromDirectory imports a class from a given directory. Proji will copy the
// structure and content of the directory and create a class based on it.
func (c *Class) ImportFromDirectory(directory string, excludeDirs []string) error {
	// Validate that the directory exists
	if !helper.DoesPathExist(directory) {
		return fmt.Errorf("given directory does not exist")
	}

	// Set class name from directory base name
	base := path.Base(directory)
	c.Name = base
	c.Label = pickLabel(c.Name)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// Skip base directory
		if directory == path {
			return nil
		}
		// Extract relative path
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		// Add file or folder to class
		if info.IsDir() {
			if helper.IsInSlice(excludeDirs, info.Name()) {
				return filepath.SkipDir
			}
			c.Folders = append(c.Folders, &Folder{Destination: relPath, Template: ""})
		} else {
			c.Files = append(c.Files, &File{Destination: relPath, Template: ""})
		}
		return nil
	})

	if err != nil {
		return err
	}

	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Directory might be empty")
	}
	return nil
}

// ImportFromURL imports a class from a given URL. The URL should point to a remote repo of one of the following code
// platforms: github, gitlab. Proji will copy the structure and content of the repo and create a class
// based on it.
func (c *Class) ImportFromURL(URL string) error {
	// Trim trailing '.git'
	if strings.HasSuffix(URL, ".git") {
		URL = URL[:len(URL)-len(".git")]
	}

	// Validate the URL
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	// Set class name from base name
	// E.g. https://github.com/nikoksr/proji -> proji is the base name
	c.Name = path.Base(u.Path)
	c.Label = pickLabel(c.Name)

	// Get repo tree (folder and file structure of remote repo)
	err = c.getRepoTree(u)
	if err != nil {
		return err
	}

	// Check if any data was loaded
	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Platform might be unsupported")
	}
	return nil
}

// getRepoTree gets the tree of the given repository and applies it to the class
func (c *Class) getRepoTree(url *url.URL) error {
	var r repo.Importer
	var err error
	escapedURL := url.Hostname() + url.EscapedPath()

	// Handle different platforms
	switch url.Hostname() {
	case "github.com":
		r, err = github.New(escapedURL)
	case "gitlab.com":
		r, err = gitlab.New(escapedURL)
	default:
		return fmt.Errorf("platform not supported yet")
	}

	if err != nil {
		return err
	}

	// Get paths and types
	paths, types, err := r.GetTreePathsAndTypes()
	if err != nil {
		return err
	}
	c.Folders, c.Files = convertPathsTypesToFoldersFiles(paths, types)

	// Set class name and label
	c.Name = r.GetRepoName()
	c.Label = pickLabel(c.Name)
	return nil
}

// convertPathsTypesToFoldersFiles converts the git types blob and tree to proji types folder and file
func convertPathsTypesToFoldersFiles(paths, types []gjson.Result) ([]*Folder, []*File) {
	// Splitting in folders and files
	folders := make([]*Folder, 0)
	files := make([]*File, 0)
	for idx, p := range paths {
		dest := p.String()

		if types[idx].String() == "tree" {
			folders = append(folders, &Folder{Destination: dest})
		} else {
			files = append(files, &File{Destination: dest})
		}
	}
	return folders, files
}

// Export exports a given class to a toml config file
func (c *Class) Export(destination string) (string, error) {
	// Create config string
	var configTxt = map[string]interface{}{
		"name":   c.Name,
		"label":  c.Label,
		"folder": c.Folders,
		"file":   c.Files,
		"script": c.Scripts,
	}

	// Export data to toml
	confName := filepath.Join(destination, "/proji-"+c.Name+".toml")
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()
	return confName, toml.NewEncoder(conf).Encode(configTxt)
}

// isEmpty checks if the class holds no data
func (c *Class) isEmpty() bool {
	if len(c.Folders) == 0 && len(c.Files) == 0 && len(c.Scripts) == 0 {
		return true
	}
	return false
}

// pickLabel dynamically picks a label based on the class name
func pickLabel(className string) string {
	nameLen := len(className)
	if nameLen < 2 {
		return strings.ToLower(className)
	}

	label := ""
	maxLabelLen := 4

	// Try to create label by separators
	seps := []string{"-", "_", ".", " "}
	parts := make([]string, 0)

	for _, d := range seps {
		parts = strings.Split(className, d)
		if len(parts) > 1 {
			break
		}
	}

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Try to create label by uppercase letters
	if !unicode.IsUpper(rune(className[0])) {
		className = string(byte(unicode.ToUpper(rune(className[0])))) + className[1:]
	}

	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	parts = re.FindAllString(className, -1)

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Pick first, mid and last byte in string
	label = string(className[0]) + string(className[nameLen/2]) + string(className[nameLen-1])
	return strings.ToLower(label)
}
