package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseConfigPath(t *testing.T) {
	bcp, err := GetBaseConfigPath()
	assert.NoError(t, err)
	assert.NotEmpty(t, bcp)
}

func TestInitConfig(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "/proji/")

	type args struct {
		path            string
		version         string
		fallbackVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All Valid 1",
			args: args{
				path:            tmpDir,
				version:         "0.20.0",
				fallbackVersion: "0.19.2",
			},
			wantErr: false,
		},
		{
			name: "All Valid 2",
			args: args{
				path:            tmpDir,
				version:         "0.19.2",
				fallbackVersion: "0.18.0",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		err := InitConfig(test.args.path, test.args.version, test.args.fallbackVersion, false)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)
		if err == nil {
			assert.DirExists(t, test.args.path, "%s\n", test.name)
		}
		_ = os.RemoveAll(test.args.path)
	}
}
