package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"io/ioutil"

	"github.com/nanobox-io/nanobox/util/config"
)

func ServiceConfigFile() string {
	switch launchSystem() {
	case "systemd":
		return "/etc/systemd/system/nanobox-openvpn.service"
	case "upstart":
		return "/etc/init/nanobox-openvpn.conf"	
	}
	return ""
}

func serviceConfig() string {
	switch launchSystem() {
	case "systemd":

		return fmt.Sprintf(`[Unit]
Description=Nanobox Openvpn Client
After=network.target

[Service]
Type=simple
EnvironmentFile=-/etc/sysconfig/network
ExecStart=%s --config %s
`, filepath.Join(config.BinDir(),BridgeClient), ConfigFile())

	case "upstart":

		return fmt.Sprintf(`
script
%s --config %s
end script`, BridgeClient, ConfigFile())

	}

	return ""
}

func CreateService() error {
	// setup config file	
	return ioutil.WriteFile(ServiceConfigFile(), []byte(serviceConfig()), 0644)
}

func StartService() error {
	switch launchSystem() {
	case "systemd":
		// systemctl start nanobox-openvpn.service
		out, err := exec.Command("systemctl", "start", "nanobox-openvpn.service").CombinedOutput() 
		if err != nil {
			return fmt.Errorf("out: %s, err: %s", out, err)
		}

	case "upstart":
		// initctl start nanobox-openvpn
		out, err := exec.Command("initctl", "start", "nanobox-openvpn").CombinedOutput() 
		if err != nil {
			return fmt.Errorf("out: %s, err: %s", out, err)
		}
	}
	return nil

}

func launchSystem() string {
	_, err := os.Stat("/sbin/systemctl")
	if err != nil {
	  return "systemd"
	}

	_, err = os.Stat("/sbin/initctl")
	if err != nil {
	  return "upstart"
	}

	return ""
}