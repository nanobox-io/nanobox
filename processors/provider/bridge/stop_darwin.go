package bridge

import (
	"fmt"
	"os/exec"
)

func stopService() error {
	out, err := exec.Command("launchctl", "stop", "io.nanobox.openvpn").CombinedOutput() 
	if err != nil {
		fmt.Errorf("out: %s, err: %s", out, err)
	}
	return nil	
}