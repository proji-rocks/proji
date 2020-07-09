//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/client9/misspell/cmd/misspell
package main

import "github.com/nikoksr/proji/cmd"

func main() {
	cmd.Execute()
}
