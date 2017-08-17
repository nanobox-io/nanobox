package hookit

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

var combined bool

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {

	// display.Streamer is an io.Writer and can be passed to DockerExec
	var stream *display.Streamer

	if !combined {
		stream = display.NewStreamer(displayLevel)
	}

	stream.CaptureOutput(true)

	out, err := util.DockerExec(container, "root", "/opt/nanobox/hooks/"+hook, []string{payload}, stream)
	if err != nil && (strings.Contains(string(out), "such file or directory") && strings.Contains(err.Error(), "bad exit code(126)")) {
		// if its a 126 the hook didnt exist
		return "", nil
	}

	outs := stream.Output()
	// the boxfile hook returns the boxfile, we shouldn't append anything.
	if hook != "boxfile" {
		if out == "" {
			out = outs
		} else if outs != "" {
			out = fmt.Sprintf("%s --- %s", out, outs)
		}
	}

	if err != nil {
		// todo: add errorfquiets for errors the hooks may stream
		if strings.Contains(outs, "INVALID BOXFILE") {
			return out, util.ErrorfQuiet("[USER] invalid node in boxfile (see output for more detail)")
		}
		return out, util.ErrorfQuiet("[HOOKS] failed to execute hook (%s) on %s: %s", hook, container, err)
	}
	return out, nil
}

func DebugExec(container, hook, payload, displayLevel string) (string, error) {
	res, err := Exec(container, hook, payload, displayLevel)

	// leave early if no error
	if err != nil {
		display.ErrorTask()
	}
	if err == nil || !registry.GetBool("debug") {
		return res, err
	}

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
