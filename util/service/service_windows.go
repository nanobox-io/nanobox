package service

import (
	"fmt"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util/config"
)

func serviceConfigFile(name string) string {
	return filepath.Join(config.BinDir(), fmt.Sprintf("%s.ini", name))
}

func startCmd(name string) []string {
	return []string{"sc", "start", name}
}

func stopCmd(name string) []string {
	return []string{"sc", "stop", name}	
}

func removeCmd(name string) []string {
	return []string{"sc", "remove", name}
}