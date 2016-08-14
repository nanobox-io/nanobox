// +build !windows,!plan9,!solaris

package provider

import (
	"os/exec"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util/config"
)

// removeDatabase ...
func (destroy Destroy) removeDatabase() error {
	return exec.Command("rm", filepath.ToSlash(filepath.Join(config.GlobalDir(), "data.db"))).Run()
}
