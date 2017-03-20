package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util/config"
)

func serviceConfigFile(name string) string {
	return filepath.Join(config.BinDir(), fmt.Sprintf("%s-config.ini", name))
}

func startCmd(name string) []string {
	return []string{"sc.exe", "start", name}
}

func running(name string) bool {
	out, err := exec.Command("sc.exe", "query", name).CombinedOutput()
	if err != nil {
		return false
	}

	if !bytes.Contains(out, []byte("RUNNING")) {
		return false
	}
	return true
}

func stopCmd(name string) []string {
	return []string{"sc.exe", "stop", name}
}

func removeCmd(name string) []string {
	return []string{"sc.exe", "delete", name}
}
