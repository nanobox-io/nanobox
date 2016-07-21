package vbox

import (
  "os/exec"
)

// Try to find VBoxManage on the path
func detectVBoxManageCmdInPath() string {
	cmd := "VBoxManage"
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	return cmd
}
