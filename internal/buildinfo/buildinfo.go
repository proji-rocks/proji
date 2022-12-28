package buildinfo

var (
	// AppName is the name of the application.
	AppName = "proji"

	// AppVersion describes the version of the app. The AppVersion gets injected during build time.
	AppVersion = ""

	// BuildDirty indicates if there were uncommitted changes in the repository when the binary was build.
	BuildDirty = false
)
