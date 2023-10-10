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
var Default = Build.All

// Aliases are mage command aliases.
var Aliases = map[string]interface{}{
	"client":    Build.Client,
	"server":    Build.Server,
	"win":       Build.WinAll,
	"winClient": Build.WinClient,
	"winServer": Build.WinServer,
	"mac":       Build.MacAll,
	"macClient": Build.MacClient,
	"macServer": Build.MacServer,

	"v-all":       Versioned.All,
	"v-server":    Versioned.Server,
	"v-client":    Versioned.Client,
	"v-win":       Versioned.WinAll,
	"v-winClient": Versioned.WinClient,
	"v-winServer": Versioned.WinServer,
	"v-mac":       Versioned.MacAll,
	"v-macClient": Versioned.MacClient,
	"v-macServer": Versioned.MacServer,
}

var envWin = map[string]string{
	"GOOS":   "windows",
	"GOARCH": "amd64",
}

var envMac = map[string]string{
	"GOOS":   "darwin",
	"GOARCH": "amd64",
}

// Build is a namespace for build related targets.
type Build mg.Namespace

// Versioned is a namespace for versioned build related targets.
type Versioned mg.Namespace

func (b Build) All() error {
	mg.Deps(Deps)
	if err := Build.Server(b); err != nil {
		return err
	}
	return Build.Client(b)
}

func (Build) Server() error {
	mg.Deps(Deps)
	fmt.Println("Building Linux server...")
	return sh.Run("go", "build", "-o", "bin/bc-server", "./cmd/bc-server")
}

func (Build) Client() error {
	mg.Deps(Deps)
	fmt.Println("Building Linux client...")
	return sh.Run("go", "build", "-o", "bin/bc-client", "./cmd/bc-client")
}

func (b Build) WinAll() error {
	mg.Deps(Deps)
	if err := Build.WinServer(b); err != nil {
		return err
	}
	return Build.WinClient(b)
}

func (Build) WinServer() error {
	mg.Deps(Deps)
	fmt.Println("Building Windows server...")
	return sh.RunWith(envWin, "go", "build", "-o", "bin/win64/bc-server.exe", "./cmd/bc-server")
}

func (Build) WinClient() error {
	mg.Deps(Deps)
	fmt.Println("Building Windows client...")
	return sh.RunWith(envWin, "go", "build", "-o", "bin/win64/bc-client.exe", "./cmd/bc-client")
}

func (b Build) MacAll() error {
	mg.Deps(Deps)
	if err := Build.MacServer(b); err != nil {
		return err
	}
	return Build.MacClient(b)
}

func (Build) MacServer() error {
	mg.Deps(Deps)
	fmt.Println("Building Mac server...")
	return sh.RunWith(envMac, "go", "build", "-o", "bin/darwin/bc-server", "./cmd/bc-server")
}

func (Build) MacClient() error {
	mg.Deps(Deps)
	fmt.Println("Building Mac client...")
	return sh.RunWith(envMac, "go", "build", "-o", "bin/darwin/bc-client", "./cmd/bc-client")
}

func (v Versioned) All(version string) error {
	mg.Deps(Deps)
	if err := Versioned.Server(v, version); err != nil {
		return err
	}
	return Versioned.Client(v, version)
}

func (Versioned) Server(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Linux server (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.Run("go", "build", "-ldflags", vf, "-o", "bin/bc-server", "./cmd/bc-server")
}

func (Versioned) Client(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Linux client (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.Run("go", "build", "-ldflags", vf, "-o", "bin/bc-client", "./cmd/bc-client")
}

func (v Versioned) WinAll(version string) error {
	mg.Deps(Deps)
	if err := Versioned.WinServer(v, version); err != nil {
		return err
	}
	return Versioned.WinClient(v, version)
}

func (Versioned) WinServer(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Windows server (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.RunWith(envWin, "go", "build", "-ldflags", vf, "-o", "bin/win64/bc-server.exe", "./cmd/bc-server")
}

func (Versioned) WinClient(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Windows client (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.RunWith(envWin, "go", "build", "-ldflags", vf, "-o", "bin/win64/bc-client.exe", "./cmd/bc-client")
}

func (v Versioned) MacAll(version string) error {
	mg.Deps(Deps)
	if err := Versioned.MacServer(v, version); err != nil {
		return err
	}
	return Versioned.MacClient(v, version)
}

func (Versioned) MacServer(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Mac server (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.RunWith(envMac, "go", "build", "-ldflags", vf, "-o", "bin/darwin/bc-server", "./cmd/bc-server")
}

func (Versioned) MacClient(version string) error {
	mg.Deps(Deps)
	fmt.Printf("Building Mac client (v%s)...\n", version)
	vf := fmt.Sprintf("-X 'github.com/milonoir/business-club-game/internal/game.Version=%s'", version)
	return sh.RunWith(envMac, "go", "build", "-ldflags", vf, "-o", "bin/darwin/bc-client", "./cmd/bc-client")
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
