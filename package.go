package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// Package map a package output by go list
// this is subset of package struct in: https://github.com/golang/go/blob/master/src/cmd/go/internal/load/pkg.go#L58
type Package struct {
	Dir        string `json:"Dir"`        // directory containing package sources
	ImportPath string `json:"ImportPath"` // import path of package in dir
	Name       string `json:"Name"`       // package name
	Target     string `json:",omitempty"` // installed target for this package (may be executable)
	Root       string `json:",omitempty"` // Go root, Go path dir, or module root dir containing this package

	Module   *ModulePublic `json:",omitempty"`         // info about package's module, if any
	Goroot   bool          `json:"Goroot,omitempty"`   // is this package in the Go root?
	Standard bool          `json:"Standard,omitempty"` // is this package part of the standard Go library?
	DepOnly  bool          `json:"DepOnly,omitempty"`  // package is only a dependency, not explicitly listed

	// Source files
	GoFiles  []string `json:"GoFiles,omitempty"`  // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles []string `json:"CgoFiles,omitempty"` // .go source files that import "C"

	// Dependency information
	Deps      []string          `json:"Deps,omitempty"` // all (recursively) imported dependencies
	Imports   []string          `json:",omitempty"`     // import paths used by this package
	ImportMap map[string]string `json:",omitempty"`     // map from source import to ImportPath (identity entries omitted)

	// Error information
	Incomplete bool            `json:"Incomplete,omitempty"` // this package or a dependency has an error
	Error      *PackageError   `json:"Error,omitempty"`      // error loading package
	DepsErrors []*PackageError `json:"DepsErrors,omitempty"` // errors loading dependencies
}

// PackageError is the error info for a package when list failed
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

// ModulePublic represents the package info of a module
type ModulePublic struct {
	Path      string        `json:",omitempty"` // module path
	Version   string        `json:",omitempty"` // module version
	Versions  []string      `json:",omitempty"` // available module versions
	Replace   *ModulePublic `json:",omitempty"` // replaced by this module
	Time      *time.Time    `json:",omitempty"` // time version was created
	Update    *ModulePublic `json:",omitempty"` // available update (with -u)
	Main      bool          `json:",omitempty"` // is this the main module?
	Indirect  bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string        `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string        `json:",omitempty"` // go version used in module
	Error     *ModuleError  `json:",omitempty"` // error loading module
}

// ModuleError represents the error loading module
type ModuleError struct {
	Err string // error text
}

func ListPackages(dir string) (map[string]*Package, error) {
	if dir == "" {
		dir = "."
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && go list -json ./...", dir))

	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("excute `go list -json ./...` command failed, err: %v, stdout: %v, stderr: %v", err, string(out), errbuf.String())
	}

	dec := json.NewDecoder(bytes.NewReader(out))
	pkgs := make(map[string]*Package, 0)

	for {
		var pkg Package
		if err := dec.Decode(&pkg); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.New("failed to read go list output")
		}
		if pkg.Error != nil {
			return nil, fmt.Errorf("failed to list package: %v", pkg.Error)
		}

		pkgs[pkg.ImportPath] = &pkg
	}
	return pkgs, nil
}
