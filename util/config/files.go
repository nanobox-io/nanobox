package config

import (
	// "fmt"
	// "io/ioutil"
	// "os"
	"path/filepath"
	// "github.com/jcelliott/lumber"
	// "github.com/nanobox-io/nanobox/util"
)

// Boxfile ...
func Boxfile() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}
