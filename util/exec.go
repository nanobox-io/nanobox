// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// Exec
// func Exec(cmd, msg string) {
// }

// SudoExec
func SudoExec(cmd, msg string) {
	fmt.Printf(stylish.Bullet(msg))

	//
	scmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v %v", os.Args[0], cmd))

	// connect standard in/outputs
	scmd.Stdin = os.Stdin
	scmd.Stdout = os.Stdout
	scmd.Stderr = os.Stderr

	// run command
	if err := scmd.Run(); err != nil {
		LogFatal("[utils/exec] scmd.Run() failed", err)
	}
}
