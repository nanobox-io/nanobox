// +build !windows

package provider

import (
	"fmt"
	"path/filepath"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
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

	// get the netfsmount options
	config, _ := models.LoadConfig()
	// additionalOptions := config.NetfsMountOpts

	// cmd = []string{"sudo", "/bin/mount", "-t", "nfs", "-o", "nolock"}
	cmd = []string{"sudo", "/bin/mount", "-t", "nfs"}
	if config.NetfsMountOpts != "" {
		cmd = append(cmd, "-o", config.NetfsMountOpts)
	}

	// the ip is hardcoded because the ip in docker is always set
	// no need to detect
	source := fmt.Sprintf("\"192.168.99.1:%s\"", local)
	cmd = append(cmd, source, host)
	if b, err := Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mount: output: %s err:%s", b, err.Error())
	}

	return nil
}
