package bridge

import (
	// "fmt"
	"os"
	"os/exec"
)

func StopService() error {
	_, err := exec.Command("sc", "stop", "nanobox-vpn").CombinedOutput()
	return err
}


func Remove() error {
	_, err := exec.Command("sc", "delete", "nanobox-vpn").CombinedOutput()
	if err != nil {
		return err
	}

	os.Remove(ServiceConfigFile())
	return nil
}