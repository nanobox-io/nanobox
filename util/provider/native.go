package provider

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// Native ...
type Native struct{}

// init ...
func init() {
	Register("native", Native{})
}

// Valid ensures docker-machine is installed and available
func (native Native) Valid() (bool, []string) {
	cmd := exec.Command("docker", "ps")

	//
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return false, []string{"docker"}
	}

	return true, nil
}

func (native Native) Status() string {
	return "Running"
}

func (native Native) BridgeRequired() bool {
	return runtime.GOOS != "linux"
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

// implode loops through the docker containers we created
// and removes each one
func (native Native) Implode() error {
	// remove any crufty containers
	cmd := exec.Command("docker", "ps", "-a")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	s := string(bytes)
	parts := strings.Split(s, "\n")
	containers := []string{}

	for _, part := range parts {
		if strings.Contains(part, "nanobox_") {
			containers = append(containers, strings.Fields(part)[0])
		}
	}

	if len(containers) == 0 {
		return nil
	}
	cmdParts := append([]string{"rm", "-f"}, containers...)
	cmd = exec.Command("docker", cmdParts...)
	cmd.Stdout = display.NewStreamer("  ")
	cmd.Stderr = display.NewStreamer("  ")

	return cmd.Run()
}

// Destroy does nothing on native
func (native Native) Destroy() error {
	// TODO: remove nanobox images

	if native.hasNetwork() {
		fmt.Print(stylish.Bullet("Removing custom docker network..."))

		cmd := exec.Command("docker", "network", "rm", "nanobox")

		cmd.Stdout = display.NewStreamer("  ")
		cmd.Stderr = display.NewStreamer("  ")

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

		config, _ := models.LoadConfig()
		ip, ipNet, err := net.ParseCIDR(config.NativeNetworkSpace)
		if err != nil {
			return err
		}

		cmd := exec.Command("docker", "network", "create", "--driver=bridge", fmt.Sprintf("--subnet=%s", ipNet.String()), "--opt=\"com.docker.network.driver.mtu=1450\"", "--opt=\"com.docker.network.bridge.name=redd0\"", fmt.Sprintf("--gateway=%s", ip.String()), "nanobox")

		cmd.Stdout = display.NewStreamer("  ")
		cmd.Stderr = display.NewStreamer("  ")

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (native Native) IsReady() bool {
	return native.hasNetwork()
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

func (native Native) ReservedIPs() (rtn []string) {
	return []string{}
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

func (native Native) SetDefaultIP(ip string) error {
	// nothing is necessary here
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

func (native Native) RequiresMount() bool {
	return false
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

//
func (native Native) RemoveEnvDir(id string) error {
	if id == "" {
		return nil
	}

	return os.RemoveAll(native.HostMntDir() + id)
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
