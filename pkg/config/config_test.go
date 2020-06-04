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
		got, err := InitConfig(test.args.path, test.args.version, false)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)
		assert.Equal(t, test.want, got, "%s\n", test.name)
		if err == nil {
			assert.DirExists(t, got, "%s\n", test.name)
		}
		_ = os.RemoveAll(test.args.path)
	}
}

func TestIsConfigUpToDate(t *testing.T) {
	type args struct {
		projiVersion  string
		configVersion string
	}
	tests := []struct {
		name       string
		args       args
		isUpToDate bool
		wantErr    bool
	}{
		{
			name: "Test IsConfigUpToDate 1",
			args: args{
				projiVersion:  "0.10.0",
				configVersion: "0.10.0",
			},
			isUpToDate: true,
			wantErr:    false,
		},
		{
			name: "Test IsConfigUpToDate 2",
			args: args{
				projiVersion:  "0.10.0",
				configVersion: "0.11.0",
			},
			isUpToDate: true,
			wantErr:    true,
		},
		{
			name: "Test IsConfigUpToDate 3",
			args: args{
				projiVersion:  "0.10.0",
				configVersion: "0.9.0",
			},
			isUpToDate: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsConfigUpToDate(tt.args.projiVersion, tt.args.configVersion)
			assert.Equal(t, tt.wantErr, err != nil, "IsConfigUpToDate() error = %v, wantErr %v", err, tt.wantErr)

			if got != tt.isUpToDate {
				t.Errorf("IsConfigUpToDate() got = %v, want %v", got, tt.isUpToDate)
			}
		})
	}
}
