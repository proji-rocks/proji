package models

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	gh "github.com/google/go-github/v31/github"
	"github.com/nikoksr/proji/pkg/config"
	"github.com/nikoksr/proji/pkg/repo"
	"github.com/nikoksr/proji/pkg/repo/github"
	"github.com/nikoksr/proji/pkg/repo/gitlab"
	"github.com/nikoksr/proji/pkg/util"
	"github.com/pelletier/go-toml"
	gl "github.com/xanzy/go-gitlab"
	"gorm.io/gorm"
)

// Class represents a proji class; the central item of proji's project creation mechanism. It holds tags for gorm and
// toml defining its storage and export/import behaviour.
type Class struct {
	ID        uint           `gorm:"primarykey" toml:"-"`
	CreatedAt time.Time      `toml:"-"`
	UpdatedAt time.Time      `toml:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_class_label,unique;" toml:"-"`
	Name      string         `gorm:"not null;size:64" toml:"name"`
	Label     string         `gorm:"index:idx_class_label,unique;not null;size:16" toml:"label"`
	Templates []*Template    `gorm:"many2many:class_templates;ForeignKey:ID;References:ID" toml:"template"`
	Plugins   []*Plugin      `gorm:"many2many:class_plugins;ForeignKey:ID;References:ID" toml:"plugin"`
	IsDefault bool           `gorm:"not null" toml:"-"`
}

// labelSeparators defines a list of rues that are used to split class names and transform them to labels.
// '%20' is for escaped paths.
var labelSeparators = []string{"-", "_", ".", " ", "%20"}

const (
	templatesKey = "templates" // Map key for template files.
	pluginsKey   = "plugins"   // Map key for plugins.
)

// NewClass returns a new class
func NewClass(name, label string, isDefault bool) *Class {
	return &Class{
		Name:      name,
		Label:     label,
		Templates: nil,
		Plugins:   nil,
		IsDefault: isDefault,
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
	file, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	err = file.Unmarshal(c)
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
	if !util.DoesPathExist(path) {
		return fmt.Errorf("given directory does not exist")
	}

	// Set class name from directory base name
	base := filepath.Base(path)
	c.Name = base
	c.Label = pickLabel(c.Name)

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
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

		// Add file or folder to class
		isFile := true
		if info.IsDir() {
			if util.IsInSlice(excludeDirs, info.Name()) {
				return filepath.SkipDir
			}
			isFile = false
		}
		c.Templates = append(c.Templates, &Template{IsFile: isFile, Path: "", Destination: relPath})
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
	c.Templates = filterAndConvertTreeEntries(importer, filters)

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
	err := util.DownloadFileIfNotExists(f, dwn)
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
	filesToDownload := make(map[string][]string)

	// All templates
	var rex *regexp.Regexp
	var templates []*Template

	for _, template := range c.Templates {
		if template.Path == "" {
			continue
		}

		// Create regex and request only once and only when necessary
		if rex == nil {
			rex = regexp.MustCompile("templates/")
			err = importer.LoadTreeEntries()
			if err != nil {
				return err
			}
			templates = filterAndConvertTreeEntries(importer, []*regexp.Regexp{rex})
		}

		if len(templates) < 1 {
			return fmt.Errorf("no templates were found in repo but class %s requires templates", c.Name)
		}

		for _, template := range templates {
			// Trim the path
			trimmedFilePath := template.Destination[len("templates/"):]
			// Add file to list only if its in the current template folder
			if strings.HasPrefix(trimmedFilePath, template.Path) {
				filesToDownload[templatesKey] = append(filesToDownload[templatesKey], trimmedFilePath)
			}
		}
	}
	for _, template := range c.Templates {
		filesToDownload[templatesKey] = append(filesToDownload[templatesKey], template.Path)
	}
	for _, plugin := range c.Plugins {
		filesToDownload[pluginsKey] = append(filesToDownload[pluginsKey], plugin.Path)
	}

	// Try and get default home dir
	var downloadDestination string
	downloadDestination, err = config.GetBaseConfigPath()
	if err != nil {
		return err
	}

	// Download scripts and templates
	// Sum of templates and scripts counts
	numFiles := len(filesToDownload[templatesKey]) + len(filesToDownload[pluginsKey])
	var wg sync.WaitGroup
	wg.Add(numFiles)
	errs := make(chan error, numFiles)

	for fileType, fileList := range filesToDownload {
		for _, file := range fileList {
			go func(fileType, file string) {
				defer wg.Done()
				src := importer.FilePathToRawURI(filepath.Join(fileType, file))
				dst := filepath.Join(downloadDestination, fileType, file)
				err = util.DownloadFileIfNotExists(dst, src)
				if err != nil {
					errs <- err
				}
			}(fileType, file)
		}
	}
	wg.Wait()
	close(errs)

	var errMsg string
	err = nil
	for e := range errs {
		if e != nil {
			errMsg += fmt.Sprintf("%s\n", e.Error())
		}
	}

	if len(errMsg) > 0 {
		err = errors.New(errMsg)
	}
	return err
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
	numFiles := len(c.Templates)
	var wg sync.WaitGroup
	wg.Add(numFiles)
	classChannel := make(chan *Class, numFiles)
	errs := make(chan error, numFiles)

	for _, template := range c.Templates {
		if !template.IsFile {
			continue
		}
		go func(template *Template) {
			defer wg.Done()
			class := NewClass("", "", false)
			packageURL, err := repo.ParseURL(URL.String() + "/" + template.Destination)
			if err != nil {
				errs <- err
				return
			}
			err = class.ImportPackage(packageURL, importer)
			if err != nil {
				errs <- err
				return
			}
			classChannel <- class
		}(template)
	}

	wg.Wait()
	close(classChannel)
	close(errs)

	for cls := range classChannel {
		if cls != nil {
			classList = append(classList, cls)
		}
	}

	err = nil
	var errMsg string
	for e := range errs {
		if e != nil {
			errMsg += fmt.Sprintf("%s\n", e.Error())
		}
	}
	if len(errMsg) > 0 {
		err = errors.New(errMsg)
	}
	return classList, err
}

// Export exports a given class to a toml config file
func (c *Class) Export(destination string) (string, error) {
	confName := filepath.Join(destination, "proji-"+c.Name+".toml")
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()
	return confName, toml.NewEncoder(conf).Order(toml.OrderPreserve).Encode(c)
}

// isEmpty checks if the class holds no data
func (c *Class) isEmpty() bool {
	if len(c.Templates) == 0 && len(c.Plugins) == 0 {
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
	parts := make([]string, 0)
	for _, d := range labelSeparators {
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

/*

	This section will be refactored asap. For now, I just want working support for packages and collections and
	don't want to spend more time refactoring it before uploading it. This works, it's just very ugly. Sorry.

*/

// GetRepoImporterFromURL returns the most suiting importer based on the code hosting platform.
func GetRepoImporterFromURL(URL *url.URL, auth *config.APIAuthentication) (repo.Importer, error) {
	var importer repo.Importer
	var err error

	switch URL.Hostname() {
	case "github.com":
		importer, err = github.New(URL, auth.GHToken)
	case "gitlab.com":
		importer, err = gitlab.New(URL, auth.GLToken)
	default:
		return nil, fmt.Errorf("platform not supported yet")
	}
	return importer, err
}

func filterAndConvertTreeEntries(importer repo.Importer, filters []*regexp.Regexp) []*Template {
	if filters == nil {
		filters = make([]*regexp.Regexp, 0)
	}

	var templates []*Template

	switch importer.(type) {
	case *github.GitHub:
		templates = filterAndConvertGHTreeEntries(importer.(*github.GitHub).TreeEntries, filters)
	case *gitlab.GitLab:
		templates = filterAndConvertGLTreeEntries(importer.(*gitlab.GitLab).TreeEntries, filters)
	default:
		return nil
	}

	return templates
}

func filterAndConvertGHTreeEntries(treeEntries []*gh.TreeEntry, filters []*regexp.Regexp) []*Template {
	if filters == nil {
		filters = make([]*regexp.Regexp, 0)
	}

	templates := make([]*Template, 0)

	for _, entry := range treeEntries {
		skip := false
		for _, filter := range filters {
			if !skip {
				skip = true
			}
			if filter.FindStringIndex(entry.GetPath()) != nil {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		isFile := false
		if entry.GetType() == "blob" {
			isFile = true
		}
		templates = append(templates, &Template{
			IsFile:      isFile,
			Path:        "",
			Destination: entry.GetPath(),
		})
	}
	return templates
}

func filterAndConvertGLTreeEntries(treeEntries []*gl.TreeNode, filters []*regexp.Regexp) []*Template {
	if filters == nil {
		filters = make([]*regexp.Regexp, 0)
	}

	templates := make([]*Template, 0)

	for _, entry := range treeEntries {
		skip := false
		for _, filter := range filters {
			if !skip {
				skip = true
			}
			if filter.FindStringIndex(entry.Path) != nil {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		isFile := false
		if entry.Type == "blob" {
			isFile = true
		}
		templates = append(templates, &Template{
			IsFile:      isFile,
			Path:        "",
			Destination: entry.Path,
		})
	}
	return templates
}
