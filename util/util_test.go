package util_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
)

// test if VboxExists works as intended
func TestVboxExists(t *testing.T) {
	cmd := "vboxmanage"

	if config.OS == "windows" {
		if installPath := os.Getenv("VBOX_MSI_INSTALL_PATH"); installPath != "" {
			cmd = filepath.Join(installPath, cmd)
		}
	}

	exists := false
	if err := exec.Command(cmd, "-v").Run(); err == nil {
		exists = true
	}

	//
	testExists := util.VboxExists()

	if testExists != exists {
		t.Error("Results don't match!")
	}
}

//
func TestMD5sMatch(t *testing.T) {
}
