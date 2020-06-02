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

	gl "github.com/xanzy/go-gitlab"

	gh "github.com/google/go-github/v31/github"

	"github.com/nikoksr/proji/pkg/config"

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

// '%20' is for escaped paths.
var labelSeparators = []string{"-", "_", ".", " ", "%20"}

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

// ImportConfig imports class data from a given config file.
func (c *Class) ImportConfig(path string) error {
	// Validate that it's a toml file
	if !strings.HasSuffix(path, ".toml") {
		return fmt.Errorf("import file has to be of type 'toml'")
	}

	// Validate config is not empty
	conf, err := os.Stat(path)
	if err != nil {
		return err
	}
	if conf.Size() == 0 {
		return fmt.Errorf("import file is empty")
	}

	// Decode the file
	_, err = toml.DecodeFile(path, &c)
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

// ImportFolderStructure imports a class from a given directory. Proji will imitate the
// structure and content of the directory and create a class based on it.
func (c *Class) ImportFolderStructure(path string, excludeDirs []string) error {
	// Validate that the directory exists
	if !helper.DoesPathExist(path) {
		return fmt.Errorf("given directory does not exist")
	}

	// Set class name from directory base name
	base := filepath.Base(path)
	c.Name = base
	c.Label = pickLabel(c.Name)

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		// Skip base directory
		if path == currentPath {
			return nil
		}
		// Extract relative path
		relPath, err := filepath.Rel(path, currentPath)
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

// ImportRepoStructure imports a class from a given URL. The URL should point to a remote repo of one of the following code
// platforms: github, gitlab. Proji will imitate the structure and content of the repo and create a class
// based on it.
func (c *Class) ImportRepoStructure(importer repo.Importer, filters []*regexp.Regexp) error {
	// Import the complete repo tree. No filters needed.
	err := importer.LoadTreeEntries()
	if err != nil {
		return err
	}
	c.Files, c.Folders = filterAndConvertTreeEntries(importer, filters)

	// Check if any data was loaded
	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Platform might be unsupported")
	}

	// Set class name from base name
	// E.g. https://github.com/nikoksr/proji -> proji is the base name
	c.Name = path.Base(importer.Repo())
	c.Label = pickLabel(c.Name)
	return nil
}

// ImportPackage imports a package from a given URL. The URL should point directly to a class config in a remote repo
// of one of the following code platforms: github, gitlab. Proji will import the class config and download its
// dependencies if necessary.
func (c *Class) ImportPackage(URL *url.URL, importer repo.Importer) error {
	// Download config
	f := filepath.Join(os.TempDir(), "/proji/configs/", filepath.Base(URL.Path))
	dwn := importer.FilePathToRawURI(filepath.Join("configs/", filepath.Base(URL.Path)))
	err := helper.DownloadFileIfNotExists(dwn, f)
	if err != nil {
		return err
	}

	// Import config
	err = c.ImportConfig(f)
	if err != nil {
		return err
	}

	// Download scripts and templates
	// Create list of necessary scripts and templates
	filesNeeded := make(map[string][]string, 0)

	// All templates
	var rex *regexp.Regexp
	var files []*File
	templatesKey := "templates"
	scriptsKey := "scripts"

	for _, folder := range c.Folders {
		if folder.Template == "" {
			continue
		}

		// Create regex and request only once and only when necessary
		if rex == nil {
			rex = regexp.MustCompile("templates/")
			err = importer.LoadTreeEntries()
			if err != nil {
				return err
			}
			files, _ = filterAndConvertTreeEntries(importer, []*regexp.Regexp{rex})
		}

		if len(files) < 1 {
			return fmt.Errorf("no templates were found in repo but class %s requires templates", c.Name)
		}

		for _, file := range files {
			// Trim the path
			trimmedFilePath := file.Destination[len("templates/"):]
			// Add file to list only if its in the current template folder
			if strings.HasPrefix(trimmedFilePath, folder.Template) {
				filesNeeded[templatesKey] = append(filesNeeded[templatesKey], trimmedFilePath)
			}
		}
	}
	for _, file := range c.Files {
		filesNeeded[templatesKey] = append(filesNeeded[templatesKey], file.Template)
	}
	for _, script := range c.Scripts {
		filesNeeded[scriptsKey] = append(filesNeeded[scriptsKey], script.Name)
	}

	// Try and get default home dir
	var downloadDestination string
	downloadDestination, err = config.GetBaseConfigPath()
	if err != nil {
		return err
	}

	// Download scripts and templates
	for fileType, fileList := range filesNeeded {
		for _, file := range fileList {
			src := importer.FilePathToRawURI(filepath.Join(fileType, file))
			dst := filepath.Join(downloadDestination, fileType, file)
			err = helper.DownloadFileIfNotExists(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportClassesFromCollection imports all classes from a given URL. A collection is a repo with multiple classes. It must include
// a folder called configs, which holds the class configs. If the classes have scripts or templates as dependencies,
// they should be put into the folders scripts/ and templates/ respectively.
func ImportClassesFromCollection(URL *url.URL, importer repo.Importer) ([]*Class, error) {
	// Get list of class configs and loop through them
	re := regexp.MustCompile(`configs/.*`)
	c := NewClass("", "", false)
	err := c.ImportRepoStructure(importer, []*regexp.Regexp{re})
	if err != nil {
		return nil, err
	}

	// Check if class is empty -> no configs found
	if c.isEmpty() {
		return nil, fmt.Errorf("no configs were found")
	}

	// Import one package at a time
	classList := make([]*Class, 0)

	for _, file := range c.Files {
		class := NewClass("", "", false)
		packageURL, err := repo.ParseURL(URL.String() + "/" + file.Destination)
		if err != nil {
			return nil, err
		}
		err = class.ImportPackage(packageURL, importer)
		if err != nil {
			return nil, err
		}
		classList = append(classList, class)
	}

	return classList, nil
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

// convertPathsNTypesToFoldersNFiles converts the git types blob and tree to proji types folder and file
func convertPathsNTypesToFoldersNFiles(paths, types []gjson.Result) ([]*Folder, []*File) {
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
	// '%20' is for escaped paths.
	seps := []string{"-", "_", ".", " ", "%20"}
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

// GetRepoImporterFromURL returns the most suiting importer based on the code hosting platform.
func GetRepoImporterFromURL(URL *url.URL) (repo.Importer, error) {
	var importer repo.Importer
	var err error

	switch URL.Hostname() {
	case "github.com":
		importer, err = github.New(URL)
	case "gitlab.com":
		importer, err = gitlab.New(URL)
	default:
		return nil, fmt.Errorf("platform not supported yet")
	}
	return importer, err
}
