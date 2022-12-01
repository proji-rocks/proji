package portability

import "github.com/cockroachdb/errors"

const (
	// FileTypeTOML is the file type toml. This is the default file type. It is used to export and import package configs.
	FileTypeTOML = "toml"
	// FileTypeJSON is the file type json. It is used to export and import package configs.
	FileTypeJSON = "json"
)

// ErrUnsupportedConfigFileType is returned if the config file is not of a supported type.
var ErrUnsupportedConfigFileType = errors.New("unsupported config file type")
