// +build !windows

package dns

// DetectHostsFile returns the location of the hosts file on unix machines
func detectHostsFile() string {
	return "/etc/hosts"
}

// detectNewlineChar returns a newline character on unix machines
func detectNewlineChar() string {
	return "\n"
}
