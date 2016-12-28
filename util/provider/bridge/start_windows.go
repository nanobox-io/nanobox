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
<<<<<<< Updated upstream:util/provider/bridge/start_windows.go
func CreateService() error {
	_, err := exec.Command("sc", "nanobox-vpn", "create", fmt.Sprintf("binpath=%s --config %s", BridgeClient, ConfigFile())).CombinedOutput()
=======
func createService() error {
	_, err := exec.Command("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf("\"%s\" --config \"%s\"", bridgeClient, configFile())).CombinedOutput()
>>>>>>> Stashed changes:processors/provider/bridge/start_windows.go
	return err
}

func StartService() error {
	_, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	return err
}