package buildinfo

import (
	"runtime/debug"
	"strings"
)

var (
	// AppName is the name of the application.
	AppName = "proji"

	// AppVersion describes the version of the app. The AppVersion gets injected during build time.
	AppVersion = ""

	// BuildDirty indicates if there were uncommitted changes in the repository when the binary was build.
	BuildDirty = false
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if AppVersion == "" {
		AppVersion = "devel"
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.modified":
			if strings.ToLower(setting.Value) == "true" {
				BuildDirty = true
			}

			if BuildDirty && AppVersion != "devel" {
				AppVersion += "+CHANGES"
			}
		}
	}
}
