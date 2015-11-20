// +build !windows

package server

import (
	"os"
	"os/signal"

	syscall "github.com/docker/docker/pkg/signal"

	"github.com/nanobox-io/nanobox/util/server/terminal"
)

func monitorTerminal(stdOutFD uintptr) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// inform the server what the starting size is
	resizeTTY(terminal.GetTTYSize(stdOutFD))

	for range sigs {
		resizeTTY(terminal.GetTTYSize(stdOutFD))
	}
}
