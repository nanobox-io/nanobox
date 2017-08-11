package config

import (
	"github.com/mitchellh/go-homedir"
	"strings"
	"testing"
)

func TestGlobalDir(t *testing.T) {
	dir := GlobalDir()
	if !strings.HasSuffix(dir, ".nanobox") {
		t.Errorf("missing nanobox suffix")
	}
}

func TestLocalDir(t *testing.T) {
	dir := LocalDir()
	// this used to be 'nanobox', i think because the boxfile was at the root level
	// since removing that, the localdir returns the current directory
	if !strings.HasSuffix(dir, "nanobox/util/config") {
		t.Errorf("local dir mismatch '%s'", dir)
	}
}

func TestLocalDirName(t *testing.T) {
	dir := LocalDirName()
	if dir != "config" {
		t.Errorf("local dir name mismatch '%s'", dir)
	}
}

func TestSSHDir(t *testing.T) {
	homedir, _ := homedir.Dir()
	if SSHDir() != homedir+"/.ssh" {
		t.Errorf("incorrect ssh directory")
	}
}
