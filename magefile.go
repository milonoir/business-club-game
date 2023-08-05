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

var envWinClient = map[string]string{
	"GOOS":   "windows",
	"GOARCH": "amd64",
}

func BuildAll() error {
	mg.Deps(Deps)
	if err := Server(); err != nil {
		return err
	}
	return Client()
}

func Server() error {
	mg.Deps(Deps)
	fmt.Println("Building server...")
	return sh.Run("go", "build", "-o", "bin/bc-server", "./cmd/bc-server")
}

func Client() error {
	mg.Deps(Deps)
	fmt.Println("Building client...")
	return sh.Run("go", "build", "-o", "bin/bc-client", "./cmd/bc-client")
}

func WinClient() error {
	mg.Deps(Deps)
	fmt.Println("Building Windows client...")
	return sh.RunWith(envWinClient, "go", "build", "-o", "bin/bc-client.exe", "./cmd/bc-client")
}

func Deps() error {
	fmt.Println("Installing deps...")
	return sh.Run("go", "mod", "download")
}

func Test() error {
	fmt.Println("Running tests...")
	result, err := sh.Output("go", "test", "-race", "-count=1", "-v", "./...")
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
