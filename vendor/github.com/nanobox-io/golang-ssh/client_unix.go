// +build !windows

package ssh

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// monWinCh watches for the system to signal a window resize and requests
// a window-change from the server.
func monWinCh(session *ssh.Session, fd uintptr) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// resize the tty if any signals received
	for range sigs {
		session.SendRequest("window-change", false, termSize(fd))
	}
}
