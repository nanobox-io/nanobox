package hookit

import (
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util"
)

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {
	return util.DockerExec(container, "/opt/nanobox/hooks/" + hook, []string{payload}, display.NewStreamer(displayLevel))
}
