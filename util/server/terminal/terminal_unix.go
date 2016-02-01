// +build !windows

package terminal

import (
	"os"
	"os/signal"

	syscall "github.com/docker/docker/pkg/signal"
)

// monitor
func monitor(stdOutFD uintptr) {

	// inform the server what the starting size is
	resizeTTY(getTTYSize(stdOutFD))

	//
	sigs := make(chan os.Signal, 1)

	//
	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// wait to resize the tty on SIGWINCH (blocking)
	for range sigs {
		resizeTTY(getTTYSize(stdOutFD))
	}
}
