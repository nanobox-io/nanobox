package bridge

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util/config"
)

func ServiceConfigFile() string {
	return filepath.Join(config.BinDir(), "bridge.ini")
}

func serviceConfig() string {
	return fmt.Sprintf(`[nanobox-vpn]
startup="%s\nanobox-vpn.exe" --config "%s"
shutdown_method=winmessage
`, config.BinDir(), ConfigFile())
}

func CreateService() error {

	// setup config file
	if err := ioutil.WriteFile(ServiceConfigFile(), []byte(serviceConfig()), 0644); err != nil {
		return err
	}

	_, err := exec.Command("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf(`%s\srvstart.exe nanobox-vpn  -c "%s"`, config.BinDir(), ServiceConfigFile())).CombinedOutput()
	return err
}

func StartService() error {
	_, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	return err
}
