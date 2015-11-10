// +build !windows

package notify

import "syscall"

// setFileLimit sets rlimit to the maximum rlimit
func setFileLimit() {
	rlm := &syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlm)
	rlm.Cur = rlm.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlm)
}
