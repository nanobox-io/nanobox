// +build !windows

package server

import (
	syscall "github.com/docker/docker/pkg/signal"
	"github.com/nanobox-io/nanobox/util/server/terminal"
	"os"
	"os/signal"
)

func monitorTerminal(stdOutFD uintptr, params string) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	for range sigs {
		w, h := terminal.GetTTYSize(stdOutFD)
		resizeTTY(stdOutFD, params, w, h)
	}
}
