package domain

import (
	"strings"
)

type generatorFn func(name string) (label string, ok bool)

const maxLabelLength = 4

func labelFromSeparators(name string) (label string, ok bool) {
	for _, separator := range []string{"-", "_", ".", " ", "%20"} {
		// Try to split the name by the separator
		parts := strings.Split(name, separator)

		// Split not successful, try next separator
		if len(parts) <= 1 {
			continue
		}

		// Split successful, try to generate a label from the parts.
		for i, part := range parts {
			if i > maxLabelLength {
				break
			}
			// Append the first letter of the part to the label
			if len(part) > 0 {
				label += string(part[0])
			}
		}

		return label, true
	}

	return "", false
}

func labelFromUpperCase(name string) (label string, ok bool) {
	// Capitalize the first letter of the name
	name = strings.ToUpper(name[:1]) + name[1:]

	// Append capital letters to the label
	for _, letter := range name {
		// Check if the letter is a capital letter
		if letter >= 'A' && letter <= 'Z' {
			label += string(letter)
		}
	}

	// If label is too short, interpret the label as invalid
	if len(label) < 2 {
		return "", false
	}

	return label, true
}

func labelFromLength(name string) (label string, ok bool) {
	nameLen := len(name)

	return string(name[0]) + string(name[nameLen/2]) + string(name[nameLen-1]), true
}

// generateLabelFromName generates a label from the given name.
func generateLabelFromName(name string) (label string) {
	defer func() {
		label = strings.ToLower(label)

		if len(label) > maxLabelLength {
			label = label[:maxLabelLength]
		}
	}()

	// Sanitize the name
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return name
	}

	// Try to generate a label from the name by splitting it by separators first, then by upper case letters, then by
	// length.
	for _, generator := range []generatorFn{
		labelFromSeparators,
		labelFromUpperCase,
		labelFromLength,
	} {
		var ok bool
		label, ok = generator(name)
		if ok {
			break
		}
	}

	return label
}
