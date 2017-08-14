package share_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"

	"github.com/nanobox-io/nanobox/util/provider/share"
)

func TestMain(m *testing.M) {
	// dont modify the actual exports
	// now we shouldnt need root :)
	provider := models.Provider{
		HostIP:  "192.168.1.2",
		MountIP: "192.168.1.4",
	}

	provider.Save()

	share.EXPORTSFILE = "/tmp/exports"
	os.Remove(share.EXPORTSFILE)
	exitCode := m.Run()
	os.Remove(share.EXPORTSFILE)
	provider.Delete()
	os.Exit(exitCode)
}

func TestShare(t *testing.T) {
	path := config.LocalDir()
	if share.Exists(path) {
		t.Errorf("the path shouldnt exist but it did")
	}

	if err := share.Add(path); err != nil {
		t.Fatal("error adding", err)
	}

	if !share.Exists(path) {
		t.Fatal("the path didnt exist when it should")
	}

	subPath := filepath.Join(path, "util")
	if share.Exists(subPath) {
		t.Fatal("subpath appears to exist when it shouldnt")
	}

	if err := share.Add(subPath); err != nil {
		t.Fatal("error adding subpath", err)
	}

	if !share.Exists(subPath) {
		t.Fatal("subpath didnt exist")
	}

	if !share.Exists(path) {
		t.Fatal("path didnt exist")
	}

	if err := share.Remove(path); err != nil {
		t.Errorf("failed to remove the path %s", err)
	}

	if share.Exists(path) {
		t.Errorf("failed to remove path")
	}

	if !share.Exists(subPath) {
		t.Errorf("removing path also removed subpath")
	}

	if err := share.Remove(subPath); err != nil {
		t.Errorf("failed to remove the subPath %s", err)
	}

	if share.Exists(path) {
		t.Errorf("path exists after it shouldnt")
	}

	if share.Exists(subPath) {
		t.Errorf("subpath exists after it shouldnt")
	}
}

// this commented version is if we pull out the recursive boxfile search from localDir()
// // TestShare tests adding parent first, so we can test if removal of parent removes child/sub path
// func TestShare(t *testing.T) {
// 	// path = /go/src/github.com/nanobox-io/nanobox/util/provider/share
// 	path := config.LocalDir()
// 	// parentPath = /go/src/github.com/nanobox-io/nanobox/util/provider
// 	parentPath := filepath.Dir(path)
// 	if share.Exists(parentPath) {
// 		t.Fatal("parent appears to exist when it shouldnt")
// 	}

// 	// add parent path to /tmp/exports
// 	if err := share.Add(parentPath); err != nil {
// 		t.Fatal("error adding parent", err)
// 	}

// 	if !share.Exists(parentPath) {
// 		t.Fatal("parent didnt exist")
// 	}

// 	if share.Exists(path) {
// 		t.Errorf("the path shouldnt exist but it did")
// 	}

// 	// add path to /tmp/exports
// 	if err := share.Add(path); err != nil {
// 		t.Fatal("error adding", err)
// 	}

// 	if !share.Exists(parentPath) {
// 		t.Fatal("parentPath should still exist")
// 	}

// 	if !share.Exists(path) {
// 		t.Fatal("the path didnt exist when it should")
// 	}

// 	// remove shares from /tmp/exports
// 	if err := share.Remove(parentPath); err != nil {
// 		t.Errorf("failed to remove the parentPath %s", err)
// 	}

// 	if !share.Exists(path) {
// 		t.Errorf("child path should exist")
// 	}

// 	if share.Exists(parentPath) {
// 		t.Errorf("parent exists after it shouldnt")
// 	}

// 	if err := share.Remove(path); err != nil {
// 		t.Errorf("failed to remove the path %s", err)
// 	}

// 	if share.Exists(path) {
// 		t.Errorf("failed to remove path")
// 	}
// }
