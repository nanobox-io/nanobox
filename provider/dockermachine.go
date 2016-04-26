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
)

type (
	DockerMachine struct {
	}
)

func init() {
	Register("docker_machine", DockerMachine{})
}

func (self DockerMachine) Valid() error {
	cmd := exec.Command("docker-machine", "version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("I could not run 'docker-machine' please make sure it is in your path")
	}
	return nil
}

func (self DockerMachine) Create() error {
	if !self.isCreated() {
		lumber.Debug("not yet created")
		// docker-machine create --driver virtualbox nanobox
		cmd := exec.Command("docker-machine", "create", "--driver", "virtualbox", "nanobox")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("create output: %s", b)
			return err
		}
	}

	return nil
}

func (self DockerMachine) Reboot() error {
	err := self.Stop()
	if err != nil {
		return err
	}
	err = self.Start()
	return err
}

func (self DockerMachine) Stop() error {
	if self.isStarted() {
		// docker-machine stop nanobox
		cmd := exec.Command("docker-machine", "stop", "nanobox")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) Destroy() error {
	if self.isCreated() {
		// docker-machine rm nanobox
		cmd := exec.Command("docker-machine", "rm", "-f", "nanobox")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) Start() error {
	if !self.isStarted() {
		// docker-machine start nanobox
		cmd := exec.Command("docker-machine", "start", "nanobox")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}

	if !self.hasNetwork() {
		lumber.Debug("not yet networked")
		// docker network create --driver=bridge --subnet=192.168.0.0/16 --opt="com.docker.network.driver.mtu=1450" --opt="com.docker.network.bridge.name=redd0" --gateway=192.168.0.1 nanobox
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "docker", "network", "create", "--driver=bridge", "--subnet=192.168.0.0/24", "--opt=\"com.docker.network.driver.mtu=1450\"", "--opt=\"com.docker.network.bridge.name=redd0\"", "--gateway=192.168.0.1", "nanobox")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("add network output: %s", b)
			return err
		}
	}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "modprobe", "ip_vs")
	_, err := cmd.CombinedOutput()

	return err
}

func (self DockerMachine) DockerEnv() error {
	// docker-machine env nanobox
	// export DOCKER_TLS_VERIFY="1"
	// export DOCKER_HOST="tcp://192.168.99.102:2376"
	// export DOCKER_CERT_PATH="/Users/lyon/.docker/machine/machines/nanobox"
	// export DOCKER_MACHINE_NAME="nanobox"
	inspect := struct {
		Driver struct {
			IPAddress string
		}
		HostOptions struct {
			EngineOptions struct {
				TlsVerify bool
			}
			AuthOptions struct {
				StorePath string
			}
		}
	}{}
	cmd := exec.Command("docker-machine", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return err
	}
	err = json.Unmarshal(b, &inspect)
	if err != nil {
		lumber.Debug("marshal: %s", b)
		return err
	}
	if inspect.HostOptions.EngineOptions.TlsVerify {
		os.Setenv("DOCKER_TLS_VERIFY", "1")
	}
	os.Setenv("DOCKER_MACHINE_NAME", "nanobox")
	os.Setenv("DOCKER_HOST", fmt.Sprintf("tcp://%s:2376", inspect.Driver.IPAddress))
	os.Setenv("DOCKER_CERT_PATH", inspect.HostOptions.AuthOptions.StorePath)
	return nil
}

func (self DockerMachine) AddIP(ip string) error {
	if !self.hasIP(ip) {
		// docker-machine ssh nanobox sudo ip addr add ${IP} dev eth1
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "add", ip, "dev", "eth1")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) RemoveIP(ip string) error {
	if self.hasIP(ip) {
		// docker-machine ssh nanobox sudo ip addr del ${IP} dev eth1
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "del", ip, "dev", "eth1")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) AddNat(ip, container_ip string) error {
	if !self.hasNatPreroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -A PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-A", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", container_ip)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	if !self.hasNatPostroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -A POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-A", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", ip)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) RemoveNat(ip, container_ip string) error {
	if self.hasNatPreroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -D PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-D", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", container_ip)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	if self.hasNatPostroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -D POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-D", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", ip)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) AddMount(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))
	if !self.hasMountLocal(local) {
		// VBoxManage sharedfolder add nanobox --name <name> --hostpath ${local} --transient
		cmd := exec.Command("VBoxManage", "sharedfolder", "add", "nanobox", "--name", name, "--hostpath", local, "--transient")
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	if !self.hasMountHost(host) {
		// create folder
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "mkdir", "-p", host)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}

		// docker-machine ssh nanobox sudo mount -t vboxsf <name> ${host}
		cmd = exec.Command("docker-machine", "ssh", "nanobox", "sudo", "mount", "-t", "vboxsf", name, host)
		b, err = cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	return nil
}

func (self DockerMachine) RemoveMount(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))
	if self.hasMountLocal(local) {
		// docker-machine ssh nanobox sudo umount ${host}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "umount", host)
		b, err := cmd.CombinedOutput()
		if err != nil {
			lumber.Debug("output: %s", b)
			return err
		}
	}
	if self.hasMountHost(host) {
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

func (self DockerMachine) isCreated() bool {
	// docker-machine status nanobox
	cmd := exec.Command("docker-machine", "status", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}
	return true
}

func (self DockerMachine) hasNetwork() bool {
	// docker-machine ssh nanobox docker network inspect nanobox
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "docker", "network", "inspect", "nanobox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("hasNetwork output: %s", b)
		return false
	}
	return true
}

func (self DockerMachine) isStarted() bool {
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

func (self DockerMachine) hasIP(ip string) bool {
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

func (self DockerMachine) hasNatPreroute(host_ip, container_ip string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "PREROUTING", "-d", host_ip, "-j", "DNAT", "--to-destination", container_ip)
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}
	return true
}

func (self DockerMachine) hasNatPostroute(host_ip, container_ip string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", host_ip)
	b, err := cmd.CombinedOutput()
	if err != nil {
		lumber.Debug("output: %s", b)
		return false
	}
	return true
}

func (self DockerMachine) hasMountHost(mount string) bool {
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

func (self DockerMachine) hasMountLocal(mount string) bool {
	// VBoxManage showvminfo nanobox --machinereadable
	cmd := exec.Command("VBoxManage", "showvminfo", "nanobox", "--machinereadable")
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
