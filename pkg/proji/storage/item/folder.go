package item

// Folder represents a class folder
type Folder struct {
	Destination string `toml:"destination"` // Where should the new folder be created?
	Template    string `toml:"template"`    // Either path of template to copy or empty when no template should be used
}
