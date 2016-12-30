package bridge

import (
	// "fmt"
	"os"
	"os/exec"
	// "strings"
)

func StopService() error {
	switch launchSystem() {
	case "systemd":
		// systemctl stop nanobox-openvpn.service
		exec.Command("systemctl", "stop", "nanobox-openvpn.service").CombinedOutput() 

		// out, err := exec.Command("systemctl", "stop", "nanobox-openvpn.service").CombinedOutput() 
		// if err != nil && !strings.Contains(err.Error(), "not loaded") {
		// 	return fmt.Errorf("out: %s, err: %s", out, err)
		// }

	case "upstart":
		// initctl stop nanobox-openvpn
		exec.Command("initctl", "stop", "nanobox-openvpn").CombinedOutput()
		// out, err := exec.Command("initctl", "stop", "nanobox-openvpn").CombinedOutput() 
		// if err != nil {
		// 	return fmt.Errorf("out: %s, err: %s", out, err)
		// }
	}

	return nil	
}

func Remove() error {
	os.Remove(ServiceConfigFile())
	return nil
}