package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-getter"
)

// DoesPathExist checks if a given path exists in the filesystem.
func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// StrToUInt converts a string into a uint.
func StrToUInt(num string) (uint, error) {
	// Parse the input
	id64, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

// WantTo waits for a valid user input to confirm if he wants to do whatever was asked for.
func WantTo(question string) bool {
	var input string
	for {
		fmt.Print(question + " [y/N] ")
		n, err := fmt.Scanln(&input)
		if err != nil {
			if err.Error() != "unexpected newline" {
				fmt.Printf("Unexpected error: %v", err)
			}
		}
		if n == 1 {
			input = strings.ToLower(input)
			if input == "n" {
				return false
			} else if input == "y" {
				return true
			}
		} else if n == 0 {
			return false
		}
	}
}

// IsInSlice returns true if a given string is found in the given slice and false if not.
func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// SkipNetworkBasedTests skips network/internet dependent tests when the env variable PROJI_SKIP_NETWORK_TESTS is set to 1
func SkipNetworkBasedTests(t *testing.T) {
	env := os.Getenv("PROJI_SKIP_NETWORK_TESTS")
	if env == "1" {
		t.Skip("Skipping network based tests")
	}
}

// CreateFolderIfNotExists creates a folder at the given path if it doesn't already exist.
func CreateFolderIfNotExists(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

// DownloadFile downloads a file from an url to the local fs.
func DownloadFile(src, dst string) error {
	return getter.GetFile(dst, src)
}

// DownloadFileIfNotExists runs downloadFile() if the destination file doesn't already exist.
func DownloadFileIfNotExists(src, dst string) error {
	_, err := os.Stat(dst)
	if os.IsNotExist(err) {
		err = DownloadFile(src, dst)
	}
	return err
}
