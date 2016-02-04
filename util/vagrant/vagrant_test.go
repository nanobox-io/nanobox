//
package vagrant_test

import (
	"os/exec"
	"testing"

  "github.com/nanobox-io/nanobox/util/vagrant"
)

// test if Exists works as intended
func TestExists(t *testing.T) {

	exists := false
	if err := exec.Command("vagrant", "-v").Run(); err == nil {
		exists = true
	}

	//
	testExists := vagrant.Exists()

	if testExists != exists {
		t.Error("Results don't match!")
	}
}
