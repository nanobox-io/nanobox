package bridge

import (
	"fmt"
	"os/exec"
)

// needed function but not used
func ServiceConfigFile() string {
	return ""
}

// create a service
func CreateService() error {
	_, err := exec.Command("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf("\"%s\" --config \"%s\"", BridgeClient, ConfigFile())).CombinedOutput()
	return err
}

func StartService() error {
	_, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	return err
}