// +build windows

package file

import (
	"fmt"
	"os/exec"
)

// Copy ...
func Copy(dst, src string) error {
	out, err := exec.Command("copy", src, dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("[util/file/file_windows] exec.Command() failed: %v - %v", err, string(out))
	}

	return nil
}
