// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
// +build !windows

//
package file

import (
	"fmt"
	"os/exec"
)

//
func Copy(src, dst string) error {
	out, err := exec.Command("cp", "-R", src+"/", dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("[util/file/file_unix] exec.Command() failed: %v - %v", err, string(out))
	}

	return nil
}
