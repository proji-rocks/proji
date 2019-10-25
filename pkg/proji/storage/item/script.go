package item

// Script represents a class script
type Script struct {
	Name       string `toml:"name"`       // Name of the script inside the scripts folder
	RunAsSudo  bool   `toml:"runAsSudo"`  // Should the script be run as root?
	ExecNumber int    `toml:"execNumber"` // When to execute the script
}
