// +build windows

package ssh

import (
	"golang.org/x/crypto/ssh"
)

// monWinCh does nothing for now on windows
func monWinCh(session *ssh.Session, fd uintptr) {
	return
}
