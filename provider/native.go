package provider

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/print"
)

type (
	Native struct {
	}
)

func init() {
	Register("native", Native{})
}

// Valid ensures docker-machine is installed and available
func (self Native) Valid() error {
	if runtime.GOOS != "linux" {
		fmt.Errorf("Native only works on linux (currently)")
	}

	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("I could not run 'docker' please make sure it is in your path")
	}
	return nil
}

// does nothing for native
func (self Native) Create() error {
	// TODO: maybe some setup stuff???
	return nil
}

// does nothing for native
func (self Native) Reboot() error {
	// TODO: nothing??

	return nil
}

// does nothing on native
func (self Native) Stop() error {
	// TODO: stop what??

	return nil
}

// does nothing on native
func (self Native) Destroy() error {
	// TODO: clean up stuff??
	if self.hasNetwork() {
		fmt.Print(stylish.Bullet("Removing custom docker network..."))

		cmd := exec.Command("docker", "network", "rm", "nanobox")
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")

		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// does nothing on native
func (self Native) Start() error {
	// TODO: some networking maybe???
	if !self.hasNetwork() {
		fmt.Print(stylish.Bullet("Setting up custom docker network..."))

		cmd := exec.Command("docker", "network", "create", "--driver=bridge", "--subnet=192.168.0.0/24", "--opt=\"com.docker.network.driver.mtu=1450\"", "--opt=\"com.docker.network.bridge.name=redd0\"", "--gateway=192.168.0.1", "nanobox")
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")

		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (self Native) HostShareDir() string {
	dir := filepath.ToSlash(filepath.Join(util.GlobalDir(), "share"))
	os.MkdirAll(dir, 0755)
	return dir
}

func (self Native) HostMntDir() string {
	dir := filepath.ToSlash(filepath.Join(util.GlobalDir(), "mnt"))
	os.MkdirAll(dir, 0755)
	return dir
}

// docker env should already be configured if docker is installed
func (self Native) DockerEnv() error {
	// ensure setup??
	return nil
}

// AddIp adds an IP into the host for host access
func (self Native) AddIP(ip string) error {
	// TODO: ???

	return nil
}

// RemoveIP removes an IP from the docker-machine vm
func (self Native) RemoveIP(ip string) error {
	// TODO: ???

	return nil
}

// AddNat adds a nat to make an container accessible to the host network stack
func (self Native) AddNat(ip, container_ip string) error {
	// TODO: ???
	return nil
}

// RemoveNat removes nat from making a container inaccessible to the host network stack
func (self Native) RemoveNat(ip, container_ip string) error {
	// TODO: ???

	return nil
}

// AddMount adds a mount into the docker-machine vm
func (self Native) AddMount(local, host string) error {
	// TODO: ???
	if !self.hasMount(host) {
		return os.Symlink(local, host)
	}
	return nil
}

func (self Native) RemoveMount(_, host string) error {
	// TODO: ???
	if self.hasMount(host) {
		return os.Remove(host)
	}
	return nil
}

func (self Native) hasNetwork() bool {
	// docker-machine ssh nanobox docker network inspect nanobox
	cmd := exec.Command("docker", "network", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("hasNetwork output: %s", b)
		return false
	}
	return true
}

func (self Native) hasMount(mount string) bool {
	fi, err := os.Stat(mount)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		lumber.Debug("Error checking mount: %s", err)
	}
	if ((fi.Mode() & os.ModeSymlink) > 0) {
		return true
	}
	return false
}