package hookit

import (
  "github.com/nanobox-io/nanobox/util/display"
)

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {
  cmd := DockerCommand(container, hook, payload)
  cmd.Stderr = display.NewStreamer(displayLevel)
  return cmd.Output()
}
