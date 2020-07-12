package packageservice

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/nikoksr/proji/pkg/domain"
)

func isPackageValid(pkg *domain.Package) error {
	if len(pkg.Name) == 0 {
		return fmt.Errorf("package needs a name")
	}
	if len(pkg.Label) == 0 {
		return fmt.Errorf("package needs a label")
	}
	if len(pkg.Templates) == 0 && len(pkg.Plugins) == 0 {
		return fmt.Errorf("package has no data")
	}
	return nil
}

// pickLabel dynamically picks a label based on the package name.
func pickLabel(packageName string) string {
	nameLen := len(packageName)
	if nameLen < 2 {
		return strings.ToLower(packageName)
	}

	label := ""
	maxLabelLen := 4

	// Try to create label by separators
	// labelSeparators defines a list of rues that are used to split package names and transform them to labels.
	// '%20' is for escaped paths.
	labelSeparators := []string{"-", "_", ".", " ", "%20"}
	parts := make([]string, 0)
	for _, d := range labelSeparators {
		parts = strings.Split(packageName, d)
		if len(parts) > 1 {
			break
		}
	}

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Try to create label by uppercase letters
	if !unicode.IsUpper(rune(packageName[0])) {
		packageName = string(byte(unicode.ToUpper(rune(packageName[0])))) + packageName[1:]
	}

	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	parts = re.FindAllString(packageName, -1)

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Pick first, mid and last byte in string
	label = string(packageName[0]) + string(packageName[nameLen/2]) + string(packageName[nameLen-1])
	return strings.ToLower(label)
}
