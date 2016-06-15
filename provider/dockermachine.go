package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/util/print"
)

// DockerMachine ...
type DockerMachine struct{}

// init ...
func init() {
	Register("docker_machine", DockerMachine{})
}

// Valid ensures docker-machine is installed and available
func (machine DockerMachine) Valid() error {

	cmd := exec.Command("docker-machine", "version")

	//
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Nanobox could not run 'docker-machine' please make sure it is in your path")
	}

	return nil
}

// Create creates the docker-machine vm
func (machine DockerMachine) Create() error {

	//
	if machine.isCreated() {
		return nil
	}

	cmd := exec.Command("docker-machine", "create", "--driver", "virtualbox", "nanobox")

	if verbose {
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")
	}

	//
	fmt.Print(stylish.ProcessStart("Starting docker-machine vm"))
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Print(stylish.ProcessEnd())

	return nil
}

// Reboot reboots the docker-machine vm
func (machine DockerMachine) Reboot() error {

	if err := machine.Stop(); err != nil {
		return err
	}

	return machine.Start()
}

// Stop stops the docker-machine vm
func (machine DockerMachine) Stop() error {

	//
	if !machine.isStarted() {
		return nil
	}

	cmd := exec.Command("docker-machine", "stop", "nanobox")

	if verbose {
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")
	}

	fmt.Print(stylish.ProcessStart("Stopping docker-machine vm"))
	if err := cmd.Run(); err != nil {
		return nil
	}
	fmt.Print(stylish.ProcessEnd())

	return nil
}

// Destroy destroys the docker-machine vm
func (machine DockerMachine) Destroy() error {

	if !machine.isCreated() {
		return nil
	}

	cmd := exec.Command("docker-machine", "rm", "-f", "nanobox")

	if verbose {
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")
	}

	fmt.Print(stylish.ProcessStart("Destroying docker-machine vm"))
	if err := cmd.Run(); err != nil {
		return nil
	}
	fmt.Print(stylish.ProcessEnd())

	return nil
}

// Start starts and bootstraps docker-machine vm
func (machine DockerMachine) Start() error {

	// start the docker-machine vm
	if !machine.isStarted() {

		cmd := exec.Command("docker-machine", "start", "nanobox")

		if verbose {
			cmd.Stdout = print.NewStreamer("  ")
			cmd.Stderr = print.NewStreamer("  ")
		}

		fmt.Print(stylish.ProcessStart("Starting docker-machine vm"))
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Print(stylish.ProcessEnd())
	}

	// create custom nanobox docker network
	if !machine.hasNetwork() {

		cmd := exec.Command("docker-machine", "ssh", "nanobox", "docker", "network", "create", "--driver=bridge", "--subnet=192.168.0.0/24", "--opt='com.docker.network.driver.mtu=1450'", "--opt='com.docker.network.bridge.name=redd0'", "--gateway=192.168.0.1", "nanobox")

		if verbose {
			cmd.Stdout = print.NewStreamer("  ")
			cmd.Stderr = print.NewStreamer("  ")
		}

		fmt.Print(stylish.Bullet("Setting up custom docker network..."))
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "modprobe", "ip_vs")

	if verbose {
		cmd.Stdout = print.NewStreamer("  ")
		cmd.Stderr = print.NewStreamer("  ")
	}

	// fmt.Print(stylish.Bullet("Ensure kernel modules are loaded..."))
	return cmd.Run()
}

// HostShareDir ...
func (machine DockerMachine) HostShareDir() string {
	return "/share/"
}

// HostMntDir ...
func (machine DockerMachine) HostMntDir() string {
	return "/mnt/sda1/"
}

// DockerEnv exports the docker connection information to the running process
func (machine DockerMachine) DockerEnv() error {
	// docker-machine env nanobox
	// export DOCKER_TLS_VERIFY="1"
	// export DOCKER_HOST="tcp://192.168.99.102:2376"
	// export DOCKER_CERT_PATH="/Users/lyon/.docker/machine/machines/nanobox"
	// export DOCKER_MACHINE_NAME="nanobox"

	// create an anonymous struct that we will populate after running inspect
	inspect := struct {
		Driver struct {
			IPAddress string
		}
		HostOptions struct {
			EngineOptions struct {
				TLSVerify bool
			}
			AuthOptions struct {
				StorePath string
			}
		}
	}{}

	// fetch the docker-machine endpoint information
	cmd := exec.Command("docker-machine", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	// marshal the json output into the anonymous struct as defined above
	err = json.Unmarshal(b, &inspect)
	if err != nil {
		lumber.Debug("marshal: %s", b)
		return err
	}

	// export TLS verify if set
	if inspect.HostOptions.EngineOptions.TLSVerify {
		os.Setenv("DOCKER_TLS_VERIFY", "1")
	}

	if inspect.Driver.IPAddress == "" {
		return fmt.Errorf("docker-machine didnt start docker properly")
	}
	// set docker environment variables for client connections
	os.Setenv("DOCKER_MACHINE_NAME", "nanobox")
	os.Setenv("DOCKER_HOST", fmt.Sprintf("tcp://%s:2376", inspect.Driver.IPAddress))
	os.Setenv("DOCKER_CERT_PATH", inspect.HostOptions.AuthOptions.StorePath)

	return nil
}

// AddIP adds an IP into the docker-machine vm for host access
func (machine DockerMachine) AddIP(ip string) error {
	if machine.hasIP(ip) {
		return nil
	}

	// docker-machine ssh nanobox sudo ip addr add ${IP} dev eth1
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "add", ip, "dev", "eth1")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	return nil
}

// RemoveIP removes an IP from the docker-machine vm
func (machine DockerMachine) RemoveIP(ip string) error {
	if !machine.hasIP(ip) {
		return nil
	}

	// docker-machine ssh nanobox sudo ip addr del ${IP} dev eth1
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "del", ip, "dev", "eth1")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	return nil
}

// AddNat adds a nat to make an container accessible to the host network stack
func (machine DockerMachine) AddNat(ip, containerIP string) error {

	// add iptables prerouting rule
	if !machine.hasNatPreroute(ip, containerIP) {
		// docker-machine ssh nanobox sudo iptables -t nat -A PREROUTING -d ${hostIP} -j DNAT --to-destination ${containerIP}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-A", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", containerIP)
		if b, err := cmd.CombinedOutput(); err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	// add iptables postrouting rule
	if !machine.hasNatPostroute(ip, containerIP) {
		// docker-machine ssh nanobox sudo iptables -t nat -A POSTROUTING -s ${containerIP} -j SNAT --to-source ${hostIP}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-A", "POSTROUTING", "-s", containerIP, "-j", "SNAT", "--to-source", ip)
		if b, err := cmd.CombinedOutput(); err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// RemoveNat removes nat from making a container inaccessible to the host network stack
func (machine DockerMachine) RemoveNat(ip, containerIP string) error {

	// remove iptables prerouting rule
	if machine.hasNatPreroute(ip, containerIP) {
		// docker-machine ssh nanobox sudo iptables -t nat -D PREROUTING -d ${hostIP} -j DNAT --to-destination ${containerIP}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-D", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", containerIP)
		if b, err := cmd.CombinedOutput(); err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	// remove iptables postrouting rule
	if machine.hasNatPostroute(ip, containerIP) {
		// docker-machine ssh nanobox sudo iptables -t nat -D POSTROUTING -s ${containerIP} -j SNAT --to-source ${hostIP}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-D", "POSTROUTING", "-s", containerIP, "-j", "SNAT", "--to-source", ip)
		if b, err := cmd.CombinedOutput(); err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// HasShare checks to see if the share exists
func (machine DockerMachine) HasShare(local, _ string) bool {
	// VBoxManage showvminfo nanobox --machinereadable
	cmd := exec.Command("VBoxManage", "showvminfo", "nanobox", "--machinereadable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match(local, output)
	if regerr != nil {
		return false
	}

	return matched
}

// AddShare adds the provided path as a shareable filesystem
func (machine DockerMachine) AddShare(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if !machine.HasShare(local, host) {
		// VBoxManage sharedfolder add nanobox --name <name> --hostpath ${local} --transient
		cmd := exec.Command("VBoxManage", "sharedfolder", "add", "nanobox", "--name", name, "--hostpath", local, "--transient")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// RemoveShare removes the provided path as a shareable filesystem; we don't care
// what the user has configured, we need to remove any shares that may have been
// setup previously
func (machine DockerMachine) RemoveShare(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if machine.HasShare(local, host) {

		// VBoxManage sharedfolder remove nanobox --name <name> --transient
		cmd := exec.Command("VBoxManage", "sharedfolder", "remove", "nanobox", "--name", name, "--transient")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// HasMount checks to see if the mount exists in the vm
func (machine DockerMachine) HasMount(mount string) bool {
	// docker-machine ssh nanobox sudo cat /proc/mounts
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "cat", "/proc/mounts")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match(mount, output)
	if regerr != nil {
		return false
	}

	return matched
}

// AddMount adds a virtualbox mount into the docker-machine vm
func (machine DockerMachine) AddMount(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if !machine.HasMount(host) {

		// create folder
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "mkdir", "-p", host)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}

		// docker-machine ssh nanobox sudo mount -t vboxsf <name> ${host}
		cmd = exec.Command("docker-machine", "ssh", "nanobox", "sudo", "mount", "-t", "vboxsf", "-o", "uid=1000,gid=1000", name, host)
		b, err = cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// RemoveMount removes a mount from the docker-machine vm
func (machine DockerMachine) RemoveMount(_, host string) error {

	if machine.HasMount(host) {
		// docker-machine ssh nanobox sudo umount ${host}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "umount", host)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	return nil
}

// HostIP inspects docker-machine to return the IP address of the vm
func (machine DockerMachine) HostIP() (string, error) {
	// create an anonymous struct that we will populate after running inspect
	inspect := struct {
		Driver struct {
			IPAddress string
		}
	}{}

	// fetch the docker-machine endpoint information
	cmd := exec.Command("docker-machine", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return "", err
	}

	// marshal the json output into the anonymous struct as defined above
	err = json.Unmarshal(b, &inspect)
	if err != nil {
		lumber.Debug("marshal: %s", b)
		return "", err
	}

	return inspect.Driver.IPAddress, nil
}

// Run a command in the vm
func (machine DockerMachine) Run(command []string) ([]byte, error) {

	// All commands need to be run in the docker machine, so we create a prefix
	context := []string{"docker-machine", "ssh", "nanobox", "sudo"}

	// now we can generate a run command combining the context with the command
	run := append(context, command...)

	// when we actually run the command, we need to pop off the first item
	cmd := exec.Command(run[0], run[1:]...)

	// run the command and return the output
	return cmd.CombinedOutput()
}

// isCreated ...
func (machine DockerMachine) isCreated() bool {
	// docker-machine status nanobox
	cmd := exec.Command("docker-machine", "status", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}
	return true
}

// hasNetwork ...
func (machine DockerMachine) hasNetwork() bool {
	// docker-machine ssh nanobox docker network inspect nanobox
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "docker", "network", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("hasNetwork output: %s", b)
		return false
	}

	return true
}

// isStarted ...
func (machine DockerMachine) isStarted() bool {
	// docker-machine status nanobox
	cmd := exec.Command("docker-machine", "status", "nanobox")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match("Running", output)
	if regerr != nil {
		return false
	}

	return matched
}

// hasIP ...
func (machine DockerMachine) hasIP(ip string) bool {
	// docker-machine ssh nanobox ip addr show dev eth1
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "ip", "addr", "show", "dev", "eth1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match(ip, output)
	if regerr != nil {
		return false
	}

	return matched
}

// hasNatPreroute ...
func (machine DockerMachine) hasNatPreroute(hostIP, containerIP string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C PREROUTING -d ${hostIP} -j DNAT --to-destination ${containerIP}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "PREROUTING", "-d", hostIP, "-j", "DNAT", "--to-destination", containerIP)
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}

	return true
}

// hasNatPostroute
func (machine DockerMachine) hasNatPostroute(hostIP, containerIP string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C POSTROUTING -s ${containerIP} -j SNAT --to-source ${hostIP}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "POSTROUTING", "-s", containerIP, "-j", "SNAT", "--to-source", hostIP)
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}

	return true
}
