package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/print"
)

// Native ...
type Native struct{}

// init ...
func init() {
	Register("native", Native{})
}

// Valid ensures docker-machine is installed and available
func (native Native) Valid() error {

	//
	if runtime.GOOS != "linux" {
		return fmt.Errorf("Native only works on linux (currently)")
	}

	cmd := exec.Command("docker", "version")

	//
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("I could not run 'docker' please make sure it is in your path")
	}

	return nil
}

// Create does nothing for native
func (native Native) Create() error {
	// TODO: maybe some setup stuff???
	return nil
}

// Reboot does nothing for native
func (native Native) Reboot() error {
	// TODO: nothing??
	return nil
}

// Stop does nothing on native
func (native Native) Stop() error {
	// TODO: stop what??
	return nil
}

// Destroy does nothing on native
func (native Native) Destroy() error {

	// TODO: clean up stuff??
	if native.hasNetwork() {
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

// Start does nothing on native
func (native Native) Start() error {

	// TODO: some networking maybe???
	if !native.hasNetwork() {
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

func (native Native) IsReady() bool {
	return true
}

// HostShareDir ...
func (native Native) HostShareDir() string {
	dir := filepath.ToSlash(filepath.Join(config.GlobalDir(), "share"))
	os.MkdirAll(dir, 0755)

	return dir + "/"
}

// HostMntDir ...
func (native Native) HostMntDir() string {
	dir := filepath.ToSlash(filepath.Join(config.GlobalDir(), "mnt"))
	os.MkdirAll(dir, 0755)

	return dir + "/"
}

// HostIP returns the loopback ip
func (native Native) HostIP() (string, error) {
	return "127.0.0.1", nil
}

// DockerEnv docker env should already be configured if docker is installed
func (native Native) DockerEnv() error {
	// ensure setup??
	return nil
}

// AddIP adds an IP into the host for host access
func (native Native) AddIP(ip string) error {
	// TODO: ???
	return nil
}

// RemoveIP removes an IP from the docker-machine vm
func (native Native) RemoveIP(ip string) error {
	// TODO: ???
	return nil
}

// AddNat adds a nat to make an container accessible to the host network stack
func (native Native) AddNat(ip, containerIP string) error {
	// TODO: ???
	return nil
}

// RemoveNat removes nat from making a container inaccessible to the host network stack
func (native Native) RemoveNat(ip, containerIP string) error {
	// TODO: ???
	return nil
}

// HasShare is not applicable for the native adapter, so will return false
func (native Native) HasShare(_, _ string) bool {
	return false
}

// AddShare is not applicable for the native adapter, so will return nil
func (native Native) AddShare(_, _ string) error {
	return nil
}

// RemoveShare is not applicable for the native adapter, so will return nil
func (native Native) RemoveShare(_, _ string) error {
	return nil
}

// HasMount will return true if the mount already exists
func (native Native) HasMount(path string) bool {
	//
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		lumber.Debug("Error checking mount: %s", err)
	}

	//
	if (fi.Mode() & os.ModeSymlink) > 0 {
		return true
	}

	return false
}

// AddMount adds a mount into the docker-machine vm
func (native Native) AddMount(local, host string) error {

	// TODO: ???
	if !native.HasMount(host) {
		if err := os.MkdirAll(filepath.Dir(host), 0755); err != nil {
			return err
		}

		return os.Symlink(local, host)
	}

	return nil
}

// RemoveMount ...
func (native Native) RemoveMount(_, host string) error {

	// TODO: ???
	if native.HasMount(host) {
		return os.Remove(host)
	}

	return nil
}

// Run will run a command on the local machine (pass-through)
func (native Native) Run(command []string) ([]byte, error) {
	// when we actually run the command, we need to pop off the first item
	cmd := exec.Command(command[0], command[1:]...)

	// run the command and return the output
	return cmd.CombinedOutput()
}

// hasNetwork ...
func (native Native) hasNetwork() bool {

	// docker-machine ssh nanobox docker network inspect nanobox
	cmd := exec.Command("docker", "network", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()

	//
	if err != nil {
		lumber.Debug("hasNetwork output: %s", b)
		return false
	}

	return true
}
