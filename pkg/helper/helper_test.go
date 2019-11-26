package helper

import (
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
