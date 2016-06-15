// Package util ...
package util

import (
	"math/rand"
	"os"
)

const (

	// VERSION is the global version for nanobox; mainly used in the update process
	// but placed here to allow access when needed (commands, processor, etc.)
	VERSION     = "1.0.0"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// IsPrivileged will return true if the current process is running under a
// privileged user, like root
func IsPrivileged() bool {

	// Execute a syscall to return the user id. If the user id is 0 then we're
	// running with root escalation.
	if os.Geteuid() != 0 {
		return true
	}

	return false
}

// RandomString ...
func RandomString(size int) string {

	//
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
