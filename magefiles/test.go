//nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test Execute unit testing.
func Test() error {
	out, err := sh.Output("go", "test", "-v", "-race", "./...", "-cover", "-coverprofile=coverage.out")
	fmt.Println(out)

	return err
}

// Cover Show HTML coverage output.
func Cover() error {
	mg.Deps(Test)
	return sh.Run("go", "tool", "cover", "-html", "coverage.out")
}
