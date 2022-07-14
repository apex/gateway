//nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// PreCommit Install pre-commit hooks
func PreCommit() {
	preCommit := sh.OutCmd("pre-commit")

	out, _ := preCommit("install", "--hook-type", "commit-msg")
	fmt.Println(out)

	out, _ = preCommit("install")
	fmt.Println(out)
}

// Install install dependencies
func Install() error {
	return sh.Run("go", "mod", "download")
}
