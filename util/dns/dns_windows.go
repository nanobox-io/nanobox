// +build windows

package dns

import (
	"os"
)

// detectHostsFile returns the location of the hosts file after expanding the
// %SystemRoot% environment variable on windows
func detectHostsFile() string {
	return os.Getenv("SystemRoot") + "\\system32\\drivers\\etc\\hosts"
}

// detectNewlineChar returns a carriage return and a newline on windows
func detectNewlineChar() string {
	return "\r\n"
}
