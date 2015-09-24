// +build linux darwin

package util

import (
	"syscall"
)

// set my rlimit to the maximum rlimit
func setFileLimit() {
	rlm := &syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlm)
	rlm.Cur = rlm.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlm)
}
