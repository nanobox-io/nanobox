package hookit

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {
	out, err := util.DockerExec(container, "root", "/opt/nanobox/hooks/"+hook, []string{payload}, display.NewStreamer(displayLevel))
	if err != nil && (strings.Contains(string(out), "No such file or directory") && strings.Contains(err.Error(), "bad exit code(126)")) {
		// if its a 126 the hook didnt exist
		return "", nil
	}
	return out, err
}

func DebugExec(container, hook, payload, displayLevel string) (string, error) {
	res, err := Exec(container, hook, payload, displayLevel)

	// leave early if no error
	if err == nil {
		return res, err
	}

	display.ErrorTask()
	err = fmt.Errorf("failed to execute %s hook: %s", hook, err.Error())
	if registry.GetBool("debug") {
		fmt.Printf("An error has occurred: \"%s\"\n", err)
		fmt.Println("Entering Debug Mode")
		fmt.Printf("  container: %s\n", container)
		fmt.Printf("  hook:      %s\n", hook)
		fmt.Printf("  payload:   %s\n", payload)
		err := console.Run(container, console.ConsoleConfig{})
		if err != nil {
			return res, fmt.Errorf("failed to establish a debug session: %s", err.Error())
		}
	}

	// try running the exec one more time.
	return Exec(container, hook, payload, displayLevel)
}
