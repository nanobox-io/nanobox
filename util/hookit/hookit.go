package hookit

import (
	"fmt"	

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/console"
)

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {
	return util.DockerExec(container, "root", "/opt/nanobox/hooks/"+hook, []string{payload}, display.NewStreamer(displayLevel))
}

func DebugExec(container, hook, payload, displayLevel string) (string, error) {
	res, err := Exec(container, hook, payload, displayLevel)
	if err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to execute %s hook: %s", hook, err.Error())
 		if registry.GetBool("debug") {
			err := console.Run(container, console.ConsoleConfig{})
			if err != nil {
				return res, fmt.Errorf("failed to establish a debug session: %s", err.Error())
			}
		}
	}

	return res, err
}