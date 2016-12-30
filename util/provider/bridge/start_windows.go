package bridge

import (
	"fmt"
	"os/exec"
	"io/ioutil"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util/update"
)


func ServiceConfigFile() string {
	return filepath.Join(publicFolder(),"bridge.ini")
}

func serviceConfig() string {
	return fmt.Sprintf(`[nanobox-vpn]
startup="%s\nanobox-vpn.exe" --config "%s"
shutdown_method=winmessage
`, publicFolder(), ConfigFile())
}

func CreateService() error {

	// setup config file
	if err := ioutil.WriteFile(ServiceConfigFile(), []byte(serviceConfig()), 0644); err != nil {
		return err
	}

	_, err := exec.Command("sc", "create", "nanobox-vpn", "binpath=", fmt.Sprintf(`%s\srvstart.exe nanobox-vpn  -c "%s"`, publicFolder(), ServiceConfigFile())).CombinedOutput()
	return err
}

func StartService() error {
	_, err := exec.Command("sc", "start", "nanobox-vpn").CombinedOutput()
	return err
}

func publicFolder() string {
	path, err := exec.LookPath(update.Name)
	fmt.Println("path", path)
	if err != nil {
		return "C:\\Program Files\\Nanobox"
	}
	return filepath.Dir(path)
}