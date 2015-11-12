//
package vagrant

import (
	"os/exec"
	"testing"
)

// test if Exists works as intended
func TestExists(t *testing.T) {

	exists := false
	if err := exec.Command("vagrant", "-v").Run(); err == nil {
		exists = true
	}

	//
	testExists := Exists()

	if testExists != exists {
		t.Error("Results don't match!")
	}
}
