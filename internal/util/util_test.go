package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoesPathExist(t *testing.T) {
	tests := []struct {
		path   string
		exists bool
	}{
		{path: "./util_test.go", exists: true},
		{path: "../../README.md", exists: true},
		{path: "../../READMENOT.md", exists: false},
		{path: "./CrYpTicFiLe.txt", exists: false},
	}

	for _, test := range tests {
		exists := DoesPathExist(test.path)
		assert.Equalf(t, test.exists, exists, "path %s expected: %v got: %v\n", test.path, test.exists, exists)
	}
}

func TestCreateFolderIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "proji-testing")

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
		err := CreateFolderIfNotExists(test.args.path)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.DirExists(t, test.args.path, "%s\n", test.name)
		}
	}

	_ = os.RemoveAll(tmpDir)
}

func TestDownloadFile(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "proji-testing")
	_ = os.RemoveAll(tmpDir)

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
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download existing file - v2",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download from private remote",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-private/master/top-secret.txt",
				dst: filepath.Join(tmpDir, "top-secret.txt"),
			},
			wantErr: true,
		},
		{
			name: "Download from invalid URL",
			args: args{
				src: "raw.githubusercontent.com/nikoksr/proji-test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "tasks.json"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		err := DownloadFile(test.args.dst, test.args.src)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	_ = os.RemoveAll(tmpDir)
}

func TestDownloadFileIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "proji-testing")
	err := CreateFolderIfNotExists(tmpDir)
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
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download a file that doesn't exist locally - v2",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download a file that already exists locally",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/CMakeLists.txt"),
			},
			wantErr: false,
		},
		{
			name: "Download from private remote",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-private/master/top-secret.txt",
				dst: filepath.Join(tmpDir, "/top-secret.txt"),
			},
			wantErr: true,
		},
		{
			name: "Download from invalid URL",
			args: args{
				src: "raw.githubusercontent.com/nikoksr/proji-test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		err := DownloadFileIfNotExists(test.args.dst, test.args.src)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	_ = os.RemoveAll(tmpDir)
}
