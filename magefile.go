//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = BuildAll

func BuildAll() error {
	mg.Deps(InstallDeps)
	if err := BuildServer(); err != nil {
		return err
	}
	return BuildClient()
}

func BuildServer() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building server...")
	return sh.Run("go", "build", "-o", "bin/bc-server", "./cmd/bc-server")
}

func BuildClient() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building client...")
	return sh.Run("go", "build", "-o", "bin/bc-client", "./cmd/bc-client")
}

func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}
