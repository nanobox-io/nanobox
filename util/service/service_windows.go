package service

import (
	"fmt"
	"path/filepath"
	"os/exec"
	"bytes"

	"github.com/nanobox-io/nanobox/util/config"
)

func serviceConfigFile(name string) string {
	return filepath.Join(config.BinDir(), fmt.Sprintf("%s-config.ini", name))
}

func startCmd(name string) []string {
	return []string{"sc", "start", name}
}

func running(name string) bool {
	out, err := exec.Command("sc", "query", name).CombinedOutput()
	if err != nil {
		return false
	}

	if !bytes.Contains(out, []byte("RUNNING")) {
		return false
	}
	return true
}

func stopCmd(name string) []string {
	return []string{"sc", "stop", name}
}

func removeCmd(name string) []string {
	return []string{"sc", "delete", name}
}
