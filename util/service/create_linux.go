package service

import (
	"fmt"
	"strings"
	"io/ioutil"
)


func Create(name string, command []string) error {
	// setup config file
	return ioutil.WriteFile(serviceConfigFile(name), []byte(serviceConfig(name, command)), 0644)
}

func serviceConfig(name string, command []string) string {
	switch launchSystem() {
	case "systemd":

		return fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
Type=simple
EnvironmentFile=-/etc/sysconfig/network
ExecStart=%s
`, name, strings.Join(command, " "))

	case "upstart":

		return fmt.Sprintf(`
script
%s
end script`, strings.Join(command, " "))

	}

	return ""
}

