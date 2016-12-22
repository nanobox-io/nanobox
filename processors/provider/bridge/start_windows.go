package bridge

import (
	"fmt"
	"os/exec"
)

// needed function but not used
func serviceConfigFile() string {
	return ""
}

// create a service
func createService() error {
	_, err := exec.Command("sc", "nanobox-vpn", "create", fmt.Sprintf("binpath=%s --config %s", bridgeClient, configFile())).CombinedOutput()
	return err
}

func startService() error {
	_, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	return err
}