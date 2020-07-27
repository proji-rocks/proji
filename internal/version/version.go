package version

import "fmt"

type version struct {
	major int
	minor int
	patch int
}

func (v version) toString() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

func Proji() string {
	return version{major: 0, minor: 20, patch: 0}.toString()
}
