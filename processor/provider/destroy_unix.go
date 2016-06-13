// +build !windows,!plan9,!solaris

package provider

import (
	"os/exec"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util"
)

// removeDatabase ...
func (providerDestroy processProviderDestroy) removeDatabase() error {
	return exec.Command("rm", filepath.ToSlash(filepath.Join(util.GlobalDir(), "data.db"))).Run()
}
