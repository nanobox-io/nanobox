package bridge

import (
	// "fmt"
	"os"
	"os/exec"
)

func StopService() error {
	exec.Command("launchctl", "stop", "io.nanobox.openvpn").CombinedOutput()

	// out, err := exec.Command("launchctl", "stop", "io.nanobox.openvpn").CombinedOutput()
	// if err != nil {
	// 	fmt.Errorf("out: %s, err: %s", out, err)
	// }
	return nil
}

func Remove() error {
	_, err := exec.Command("launchctl", "remove", "io.nanobox.openvpn").CombinedOutput()
	if err != nil {
		return err
	}

	os.Remove(ServiceConfigFile())
	return nil
}
