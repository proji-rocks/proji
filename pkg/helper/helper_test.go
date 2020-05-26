package helper

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoesPathExist(t *testing.T) {
	tests := []struct {
		path   string
		exists bool
	}{
		{path: "./helper_test.go", exists: true},
		{path: "../../README.md", exists: true},
		{path: "../../READMENOT.md", exists: false},
		{path: "./CrYpTicFiLe.txt", exists: false},
	}

	for _, test := range tests {
		exists := DoesPathExist(test.path)
		assert.Equal(t, test.exists, exists)
	}
}

func TestStrToUInt(t *testing.T) {
	tests := []struct {
		numAsStr string
		expNum   uint
		err      error
	}{
		{numAsStr: "0", expNum: 0, err: nil},
		{numAsStr: "2142534513", expNum: 2142534513, err: nil},
		{numAsStr: "-1", expNum: 0, err: &strconv.NumError{}},
		{numAsStr: "1231231233123123123123231", expNum: 0, err: &strconv.NumError{}},
	}

	for _, test := range tests {
		actNum, err := StrToUInt(test.numAsStr)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expNum, actNum)
	}
}

func TestIsInSlice(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				slice: []string{"test1", "test2", "test3", "test4", "test5"},
				val:   "test5",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				slice: []string{"test1", "test2", "test3", "test4", "test5"},
				val:   "test000",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				slice: make([]string, 0),
				val:   "test",
			},
			want: false,
		},
	}
	for _, test := range tests {
		got := IsInSlice(test.args.slice, test.args.val)
		assert.Equal(t, test.want, got)
	}
}

func Test_createFolderIfNotExists(t *testing.T) {
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

func Test_downloadFile(t *testing.T) {
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
			name: "Download from private repo",
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
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: true,
		},
		{
			name: "Download file to invalid location",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/this/path/does/not/exist/CMakeLists.txt"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		err := DownloadFile(test.args.src, test.args.dst)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	_ = os.RemoveAll(tmpDir)
}

func Test_downloadFileIfNotExists(t *testing.T) {
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
				src: "raw.githubusercontent.com/nikoksr/proji-test/master/.vscode/tasks.json",
				dst: filepath.Join(tmpDir, "/tasks.json"),
			},
			wantErr: false,
		},
		{
			name: "Download file to invalid location",
			args: args{
				src: "https://raw.githubusercontent.com/nikoksr/proji-test/master/CMakeLists.txt",
				dst: filepath.Join(tmpDir, "/this/path/does/not/exist/CMakeLists.txt"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		err := DownloadFileIfNotExists(test.args.src, test.args.dst)
		assert.Equal(t, test.wantErr, err != nil, "%s\n", test.name)

		if !test.wantErr {
			assert.FileExists(t, test.args.dst, "%s\n", test.name)
		}
	}

	_ = os.RemoveAll(tmpDir)
}
