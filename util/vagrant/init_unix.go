// +build !windows

package vagrant

func sshLocation() string {
	return "~/.ssh"
}