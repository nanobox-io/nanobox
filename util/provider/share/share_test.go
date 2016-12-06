package share_test

import (
	"os"
	"os/exec"
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

	exec.Command("touch", "/tmp/exports").Run()
	share.EXPORTSFILE = "/tmp/exports"
	exitCode := m.Run()
	os.Remove("/tmp/exports")
	provider.Destroy()
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
