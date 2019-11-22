package item

// Script represents a class script
type Script struct {
	Name       string   `toml:"name"`       // Name of the script inside the scripts folder
	Type       string   `toml:"type"`       // Type of the script (pre or post) determines if the script should be run before or after project creation
	ExecNumber int      `toml:"execNumber"` // When to execute the script
	RunAsSudo  bool     `toml:"runAsSudo"`  // Should the script be run as root?
	Args       []string `toml:"args"`       // Arguments for the script that would normally be passed via CLI
}
