package bridge

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jcelliott/lumber"

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

	// make sure we actually have to do this part
	if out, _ := exec.Command("sc", "query", "nanobox-vpn").CombinedOutput(); !strings.Contains(string(out), "service does not exist") {
		return nil
	}

	// setup config file
	if err := ioutil.WriteFile(ServiceConfigFile(), []byte(serviceConfig()), 0644); err != nil {
		return err
	}

	// the service may have been created this should clean out any old version
	// we arent catching errors just incase they dont exist
	StopService()
	Remove()
	out, err := exec.Command("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf(`%s\srvstart.exe nanobox-vpn  -c "%s"`, config.BinDir(), ServiceConfigFile())).CombinedOutput()
	lumber.Info("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf(`%s\srvstart.exe nanobox-vpn  -c "%s"`, config.BinDir(), ServiceConfigFile()))
	lumber.Info("\n\nout: %s\n\n", out)
	if err != nil {
		return fmt.Errorf("%s: %s", out, err)
	}
	fmt.Printf("\n\nout: %s\n\n", out)

	return err
}

func StartService() error {
	out, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", out, err)
	}
	return err
}
