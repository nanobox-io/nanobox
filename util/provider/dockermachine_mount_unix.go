// +build !windows

package provider

import (
	"fmt"
	"path/filepath"

	"github.com/jcelliott/lumber"
)

// Mount mounts a share on a guest machine
func (machine DockerMachine) addNetfsMount(local, host string) error {
	// make local the actual path instead of the link
	local, _ = filepath.EvalSymlinks(local)

	// ensure nfs-client is running
	cmd := []string{"sudo", "/usr/local/etc/init.d/nfs-client", "start"}
	if b, err := Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("nfs-client:%s", err.Error())
	}

	// ensure the destination directory exists
	cmd = []string{"sudo", "/bin/mkdir", "-p", host}
	if b, err := Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mkdir:%s", err.Error())
	}

	// TODO: this IP shouldn't be hardcoded, needs to be figured out!
	source := fmt.Sprintf("\"192.168.99.1:%s\"", local)
	cmd = []string{"sudo", "/bin/mount", "-t", "nfs", source, host}
	if b, err := Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mount: output: %s err:%s", b, err.Error())
	}

	return nil

}
