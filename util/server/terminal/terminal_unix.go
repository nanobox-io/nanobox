// +build !windows

package terminal

import (
	"os"
	"os/signal"

	syscall "github.com/docker/docker/pkg/signal"
)

// monitor
func monitor(stdOutFD uintptr) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// inform the server what the starting size is
	resizeTTY(getTTYSize(stdOutFD))

	// resize the tty for any signals received
	for range sigs {
		resizeTTY(getTTYSize(stdOutFD))
	}
}
