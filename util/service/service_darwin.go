package service

import (
	"fmt"
	"os/exec"
	"bytes"
)

func serviceConfigFile(name string) string {
	return fmt.Sprintf("/Library/LaunchDaemons/io.%s.plist", name)
}

func startCmd(name string) []string {
	return []string{"launchctl", "start", fmt.Sprintf("io.%s", name)}
}

func Running(name string) bool {
	out, err := exec.Command("launchctl", "list", name).CombinedOutput()
	if err != nil {
		return false
	}

	if !bytes.Contains(out, []byte("PID")) {
		return false
	}
	return true
}

func stopCmd(name string) []string {
	return []string{"launchctl", "stop", fmt.Sprintf("io.%s", name)}
}

func removeCmd(name string) []string {
	return []string{"launchctl", "remove", fmt.Sprintf("io.%s", name)}
}
