package item

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly"
	"github.com/nikoksr/proji/pkg/helper"

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
		return fmt.Errorf("Import file has to be of type 'toml'")
	}

	// Validate config is not empty
	conf, err := os.Stat(configName)
	if err != nil {
		return err
	}
	if conf.Size() == 0 {
		return fmt.Errorf("Import file is empty")
	}

	// Decode the file
	_, err = toml.DecodeFile(configName, &c)
	if err != nil {
		return err
	}

	if len(c.Name) < 1 {
		return fmt.Errorf("Name cannot be an empty string")
	}
	if len(c.Label) < 1 {
		return fmt.Errorf("Label cannot be an empty string")
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
		return fmt.Errorf("Given directory does not exist")
	}

	// Set class name from directory base name
	base := path.Base(directory)
	c.Name = base
	c.Label = pickLabel(c.Name)

	// This map of directories that should be skipped might be moved to the main config
	// file so that it's editable and extensible.
	excludeDirs = append(excludeDirs, []string{".git", ".env"}...)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// Skip base directory
		if directory == path {
			return nil
		}
		// Extract relative path
		relPath, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		// Add file or folder to class
		if info.IsDir() {
			c.Folders = append(c.Folders, &Folder{Destination: relPath, Template: ""})
			if helper.IsInSlice(excludeDirs, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			c.Files = append(c.Files, &File{Destination: relPath, Template: ""})
		}
		return nil
	})

	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Directory might be empty")
	}
	return err
}

// ImportFromURL imports a class from a given URL. The URL should point to a remote repo of one of the following code
// platforms: github, gitlab, bitbucket. Proji will copy the structure and content of the repo and create a class
// based on it.
func (c *Class) ImportFromURL(URL string, excludes []string) error {
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
	base := path.Base(u.Path)
	c.Name = base
	c.Label = pickLabel(c.Name)

	// List of directories and files that should be skipped.
	excludes = append(excludes, []string{".git", ".env"}...)

	repoStructure := crawlRepo(URL, excludes)
	c.Folders, c.Files = cleanUpURLList(repoStructure, URL)

	if c.isEmpty() {
		return fmt.Errorf("no relevant data was found. Website might be unsupported")
	}
	return nil
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
	confName := destination + "/proji-" + c.Name + ".toml"
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()
	return confName, toml.NewEncoder(conf).Encode(configTxt)
}

func (c *Class) isEmpty() bool {
	if len(c.Folders) == 0 && len(c.Files) == 0 && len(c.Scripts) == 0 {
		return true
	}
	return false
}

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

func crawlRepo(URL string, excludes []string) []string {
	// Parse excludes to regex slice
	excRegex := make([]*regexp.Regexp, 0)
	for _, exc := range excludes {
		excRegex = append(excRegex, regexp.MustCompile(URL+"/(?:blob|tree)/master/"+exc))
	}

	// Will hold the repo structure - files and folders.
	repo := make([]string, 0)

	var c = colly.NewCollector(
		colly.URLFilters(regexp.MustCompile(URL+"/(?:blob|tree)/master/.*")),
		colly.DisallowedURLFilters(excRegex...),
		colly.Async(true),
	)

	_ = c.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 5 * time.Second,
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		_ = c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		repo = append(repo, r.URL.String())
	})

	_ = c.Visit(URL + "/blob/master/")
	_ = c.Visit(URL + "/tree/master/")
	c.Wait()
	return repo
}

func cleanUpURLList(repoStructure []string, baseURL string) ([]*Folder, []*File) {
	// Elements 0 and 1 are the base URLs
	// Cut out the two base URLs
	repoStructure = repoStructure[2:]

	// Pre-sort
	sort.Strings(repoStructure)

	// Make paths relative
	baseTree := baseURL + "/tree/master/"
	lenBaseTree := len(baseTree)
	baseBlob := baseURL + "/blob/master/"
	lenBaseBlob := len(baseBlob)

	folders := make([]*Folder, 0)
	files := make([]*File, 0)

	for _, URL := range repoStructure {
		if strings.HasPrefix(URL, baseTree) {
			folders = append(folders, &Folder{URL[lenBaseTree:], ""})
		} else {
			files = append(files, &File{URL[lenBaseBlob:], ""})
		}
	}
	return folders, files
}
