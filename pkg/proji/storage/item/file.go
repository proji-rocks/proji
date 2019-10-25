package item

// File represents a class file
type File struct {
	Destination string `toml:"destination"` // Where should the new file be created?
	Template    string `toml:"template"`    // Either path of template to copy or empty when no template should be used
}
