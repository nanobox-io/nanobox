package hookit

import (
	"fmt"
	"strings"
	"io"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

var combined bool

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {

	var stream io.Writer

	if !combined {
		stream = display.NewStreamer(displayLevel)
	}

	out, err := util.DockerExec(container, "root", "/opt/nanobox/hooks/"+hook, []string{payload}, stream)
	if err != nil && (strings.Contains(string(out), "such file or directory") && strings.Contains(err.Error(), "bad exit code(126)")) {
		// if its a 126 the hook didnt exist
		return "", nil
	}

	if err != nil {
		return out, util.ErrorfQuiet("failed to execute hook (%s) on %s: %s", hook, container, err)
	}
	return out, nil
}

func DebugExec(container, hook, payload, displayLevel string) (string, error) {
	res, err := Exec(container, hook, payload, displayLevel)

	// leave early if no error
	if err == nil || !registry.GetBool("debug") {
		return res, err
	}
	display.ErrorTask()

	combined = true
	res, err = Exec(container, hook, payload, displayLevel)

	fmt.Printf("failed to execute %s hook: %s\n", hook, err)
	fmt.Printf("An error has occurred: \"%s\"\n", res)
	fmt.Println("Entering Debug Mode")
	fmt.Printf("  container: %s\n", container)
	fmt.Printf("  hook:      %s\n", hook)
	fmt.Printf("  payload:   %s\n", payload)
	err = console.Run(container, console.ConsoleConfig{})
	if err != nil {
		return res, fmt.Errorf("failed to establish a debug session: %s", err.Error())
	}
	combined = false

	// try running the exec one more time.
	return Exec(container, hook, payload, displayLevel)
}
