package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "", want: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBaseConfigPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseConfigPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBaseConfigPath() got = %v, want %v", got, tt.want)
			}
		})
	}
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
			name: "All Valid",
			args: args{
				path:    tmpDir,
				version: "0.18.1",
			},
			want:    tmpDir,
			wantErr: false,
		},
		{
			name: "Invalid Version",
			args: args{
				path:    tmpDir,
				version: "__NotAVersionNumber!!",
			},
			want:    "",
			wantErr: true,
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

func Test_createFolderIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "/proji/")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create folder that doesn't exist - v1",
			args:    args{path: filepath.Join(tmpDir, "/1BnnhHku32kjchWl")},
			wantErr: false,
		},
		{
			name:    "Create folder that doesn't exist - v2",
			args:    args{path: filepath.Join(tmpDir, "/nsaido398Ncaa34J")},
			wantErr: false,
		},
		{
			name:    "Create folder that does exist - v1",
			args:    args{path: filepath.Join(tmpDir, "/1BnnhHku32kjchWl")},
			wantErr: false,
		},
		{
			name:    "Create folder that does exist - v2",
			args:    args{path: tmpDir},
			wantErr: false,
		},
	}

	for _, test := range tests {
		err := createFolderIfNotExists(test.args.path)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.DirExists(t, test.args.path, "%s\n", test.name)
		}
	}

	err := os.RemoveAll(tmpDir)
	assert.NoError(t, err)
}

func Test_downloadFile(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "/proji/")
	err := createFolderIfNotExists(tmpDir)
	assert.NoError(t, err)

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Download existing file - v1",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download existing file - v2",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download from private repo",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-private/master/top-secret.txt",
				dst: filepath.Join(tmpDir, "/top-secret.txt"),
			},
			wantErr: true,
		},
		{
			name: "Download from invalid URL",
			args: args{
				src: "raw.githubusercontent.com/nikoksr/proji_test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: true,
		},
		{
			name: "Download file to invalid location",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/this/path/does/not/exist/CMakeLists.txt"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		err := downloadFile(test.args.src, test.args.dst)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	err = os.RemoveAll(tmpDir)
	assert.NoError(t, err)
}

func Test_downloadFileIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "/proji/")
	err := createFolderIfNotExists(tmpDir)
	assert.NoError(t, err)

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Download a file that doesn't exist locally - v1",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download a file that doesn't exist locally - v2",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download a file that already exists locally",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download from private repo",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-private/master/top-secret.txt",
				dst: filepath.Join(tmpDir, "/top-secret.txt"),
			},
			wantErr: true,
		},
		{
			name: "Download from invalid URL",
			args: args{
				src: "raw.githubusercontent.com/nikoksr/proji_test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download file to invalid location",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji_test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/this/path/does/not/exist/CMakeLists.txt"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		err := downloadFileIfNotExists(test.args.src, test.args.dst)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	err = os.RemoveAll(tmpDir)
	assert.NoError(t, err)
}
