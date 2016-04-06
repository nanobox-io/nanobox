package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"os/exec"
	"regexp"
)

type (
	DockerMachine struct {
	}
)

func init() {
	Register("docker_machine", DockerMachine{})
}

func (self DockerMachine) isCreated() bool {
	// docker-machine status nanobox
	cmd := exec.Command("docker-machine", "status", "nanobox")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func (selfDockerMachine) hasNetwork() bool {
	// docker-machine ssh nanobox docker network inspect nanobox
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "docker", "network", "inspect", "nanobox")
	err := cmd.Run()
	if err != nil {
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
	return mached
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
	return mached
}

func (self DockerMachine) hasNatPreroute(host_ip, container_ip string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "PREROUTING", "-d", host_ip, "-j", "DNAT", "--to-destination", container_ip)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func (self DockerMachine) hasNatPostroute(host_ip, container_ip string) bool {
	// docker-machine ssh nanobox sudo iptables -t nat -C POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
	cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "/usr/local/sbin/iptables", "-t", "nat", "-C", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", host_ip)
	err := cmd.Run()
	if err != nil {
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
	return mached
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
	return mached
}

func (self DockerMachine) Create() error {
	if !self.isCreated() {
		// docker-machine create --driver virtualbox nanobox
		cmd := exec.Command("docker-machine", "create", "--driver", "virtualbox", "nanobox")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if !self.hasNetwork() {
		// docker network create --driver=bridge --subnet=192.168.0.0/16 --opt="com.docker.network.driver.mtu=1450" --opt="com.docker.network.bridge.name=redd0" --gateway=192.168.0.1 nanobox
		cmd := exec.Command("docker", "network", "create", "--driver=bridge", "--subnet=192.168.0.0/16", "--opt=\"com.docker.network.driver.mtu=1450\"", "--opt=\"com.docker.network.bridge.name=redd0\"", "--gateway=192.168.0.1", "nanobox")
		err := cmd.Run()
		if err != nil {
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
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) Destroy() error {
	if self.isCreated() {
		// docker-machine rm nanobox
		cmd := exec.Command("docker-machine", "rm", "nanobox")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) Start() error {
	if !self.isStarted() {
		// docker-machine start nanobox
		cmd := exec.Command("docker-machine", "start", "nanobox")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) AddIP(ip string) error {
	if !self.hasIP(ip) {
		// docker-machine ssh nanobox sudo ip addr add ${IP} dev eth1
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "add", ip, "dev", "eth1")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) RemoveIP(ip string) error {
	if self.hasIP(ip) {
		// docker-machine ssh nanobox sudo ip addr del ${IP} dev eth1
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "ip", "addr", "del", ip, "dev", "eth1")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) AddNat(ip, container_ip string) error {
	if !self.hasNatPreroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -A PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "iptables", "-t", "nat", "-A", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", container_ip)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if !self.hasNatPostroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -A POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "iptables", "-t", "nat", "-A", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", ip)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self DockerMachine) RemoveNat(ip, container_ip string) error {
	if self.hasNatPreroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -D PREROUTING -d ${host_ip} -j DNAT --to-destination ${container_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "iptables", "-t", "nat", "-D", "PREROUTING", "-d", ip, "-j", "DNAT", "--to-destination", container_ip)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if self.hasNatPostroute(ip, container_ip) {
		// docker-machine ssh nanobox sudo iptables -t nat -D POSTROUTING -s ${container_ip} -j SNAT --to-source ${host_ip}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "iptables", "-t", "nat", "-D", "POSTROUTING", "-s", container_ip, "-j", "SNAT", "--to-source", ip)
		err := cmd.Run()
		if err != nil {
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
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if !self.hasMountHost(host) {
		// docker-machine ssh nanobox sudo mount -t vboxsf <name> ${host}
		cmd := exec.Command("docker-machine", "ssh", "nanobox", "sudo", "mount", "-t", "vboxsf", name, host)
		err := cmd.Run()
		if err != nil {
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
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if self.hasMountHost(host) {
		// VBoxManage sharedfolder remove nanobox --name <name> --transient
		cmd := exec.Command("VBoxManage", "sharedfolder", "remove", "nanobox", "--name", name, "--transient")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
