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
	if !strings.HasSuffix(dir, "nanobox") {
		t.Errorf("missing nanobox suffix")
	}
}

func TestLocalDirName(t *testing.T) {
	if LocalDirName() != "nanobox" {
		t.Errorf("local dir name mismatch")
	}
}

func TestBinDir(t *testing.T) {
	dir := BinDir()
	if !strings.HasSuffix(dir, ".nanobox/bin") {
		t.Errorf("bin dir failure")
	}
}

func TestSSHDir(t *testing.T) {
	homedir, _ := homedir.Dir()
	if SSHDir() != homedir+"/.ssh" {
		t.Errorf("incorrect ssh directory")
	}
}
