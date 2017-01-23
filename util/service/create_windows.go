package service

import (
	"fmt"
	"os/exec"
	"strings"
	"io/ioutil"

	"github.com/nanobox-io/nanobox/util/config"

)

func Create(name string, command []string) error {

	// make sure we actually have to do this part
	if out, _ := exec.Command("sc", "query", name).CombinedOutput(); !strings.Contains(string(out), "service does not exist") {
		return nil
	}

	// setup config file
	if err := ioutil.WriteFile(serviceConfigFile(name), []byte(serviceConfig(name, command)), 0644); err != nil {
		return err
	}

	// the service may have been created this should clean out any old version
	// we arent catching errors just incase they dont exist
	Stop(name)
	Remove(name)
	
	out, err := exec.Command("sc", "create", name, "binpath=", fmt.Sprintf(`%s\srvstart.exe %s -c "%s"`, config.BinDir(), name, serviceConfigFile(name))).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", out, err)
	}
	fmt.Printf("\n\nout: %s\n\n", out)

	return err
}

func serviceConfig(name string, command []string) string {
	return fmt.Sprintf(`[%s]
startup=%s
shutdown_method=winmessage
`, name, strings.Join(command, " "))
}
