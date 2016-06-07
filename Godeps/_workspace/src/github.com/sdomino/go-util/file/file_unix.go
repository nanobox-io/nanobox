// +build !windows

package file

import (
	"fmt"
	"os/exec"
)

// Copy ...
func Copy(src, dst string) error {
	out, err := exec.Command("cp", "-R", src, dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("[util/file/file_unix] exec.Command() failed: %v - %v", err, string(out))
	}

	return nil
}
