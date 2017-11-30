package util

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
)

func OsDetect() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getDarwin()
	case "windows":
		return "windows", nil
	case "linux":
		return "linux", nil
	}

	return "", fmt.Errorf("Unsupported operating system. Please contact support.")
}

func getDarwin() (string, error) {
	out, err := exec.Command("/usr/bin/sw_vers", "-productVersion").Output()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve version - %s", err.Error())
	}
	r, _ := regexp.Compile("10\\.([0-9]+).*")
	match := r.FindStringSubmatch(string(out))
	if len(match) != 2 {
		return "", fmt.Errorf("Failed to parse version")
	}

	return toDarwin(match[1])
}

func toDarwin(v string) (string, error) {
	switch v {
	case "12":
		return "sierra", nil
	case "13":
		return "high sierra", nil
	default:
		return "incompatible", fmt.Errorf("Incompatible OSX version. Please contact support.")
	}
}
