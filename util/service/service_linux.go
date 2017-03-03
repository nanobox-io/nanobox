package service

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
)

func serviceConfigFile(name string) string {
	fmtString := ""
	switch launchSystem() {
	case "systemd":
		fmtString = "/etc/systemd/system/%s.service"
	case "upstart":
		fmtString = "/etc/init/%s.conf"
	}
	return fmt.Sprintf(fmtString, name)
}

func launchSystem() string {
	_, err := os.Stat("/sbin/systemctl")
	if err == nil {
		return "systemd"
	}

	_, err = os.Stat("/usr/bin/systemctl")
	if err == nil {
		return "systemd"
	}

	_, err = os.Stat("/sbin/initctl")
	if err == nil {
		return "upstart"
	}

	return ""
}

func startCmd(name string) []string {
	switch launchSystem() {
	case "systemd":
		// systemctl start nanobox-openvpn.service
		return []string{"systemctl", "start", fmt.Sprintf("%s.service", name)}
	case "upstart":
		// initctl start nanobox-openvpn
		return []string{"initctl", "start", name}
	}

	return nil
}

func running(name string) bool {
	switch launchSystem() {
	case "systemd":
		out, err := exec.Command("systemctl", "status", name).CombinedOutput()
		if err != nil {
			return false
		}

		if !bytes.Contains(out, []byte("running")) {
			return false
		}
		return true
	case "upstart":
		out, err := exec.Command("initctl", "status", name).CombinedOutput()
		if err != nil {
			return false
		}

		if !bytes.Contains(out, []byte("running")) {
			return false
		}
		return true
	}

	return false
}

func stopCmd(name string) []string {
	switch launchSystem() {
	case "systemd":
		// systemctl start nanobox-openvpn.service
		return []string{"systemctl", "stop", fmt.Sprintf("%s.service", name)}
	case "upstart":
		// initctl start nanobox-openvpn
		return []string{"initctl", "stop", name}
	}

	return nil
}

func removeCmd(name string) []string {
	return nil
}
