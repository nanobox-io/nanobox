package bridge

import (
	"os/exec"
	"github.com/nanobox-io/nanobox/commands/server"
)

type Bridge struct{}

// not being used yet. but could be
type Response struct {
	Output   string
	ExitCode int
}

var runningBridge *exec.Cmd

func init() {
	server.Register(&Bridge{})
}
