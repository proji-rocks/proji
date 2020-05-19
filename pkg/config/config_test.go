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
		path    string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "All Valid 1",
			args: args{
				path:    tmpDir,
				version: "0.19.0",
			},
			want:    tmpDir,
			wantErr: false,
		},
		{
			name: "All Valid 2",
			args: args{
				path:    tmpDir,
				version: "0.18.1",
			},
			want:    tmpDir,
			wantErr: false,
		},
	}

	for _, test := range tests {
		got, err := InitConfig(test.args.path, test.args.version)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)
		assert.Equal(t, test.want, got, "%s\n", test.name)
		if err == nil {
			assert.DirExists(t, got, "%s\n", test.name)
		}
		err = os.RemoveAll(test.args.path)
		assert.NoError(t, err, "%s\n", test.name)
	}
}
