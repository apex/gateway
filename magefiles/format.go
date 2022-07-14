//nolint
package main

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
)

// Lint Runs golangci-lint checks over the code.
func Lint() error {
	out, err := sh.Output(
		"golangci-lint",
		"run", "./...",
		"--allow-parallel-runners",
		"--skip-dirs", `(node_modules|magefiles|\.serverless|mod|bin|vendor|\.github|\.git)`,
	)

	fmt.Println(out)
	return err
}

// Format Runs gofmt over the code.
func Format() error {
	list, err := sh.Output("go", "list", "-f", "{{.Dir}}", "./...")
	if err != nil {
		return err
	}

	args := []string{
		"-s", "-w", "-l",
	}

	for _, dir := range strings.Split(list, "\n") {
		if dir == "" {
			continue
		}

		args = append(args, dir)
	}

	out, err := sh.Output("gofmt", args...)
	fmt.Println(out)

	return err
}
