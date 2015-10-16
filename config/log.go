// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
)

//
var (
	Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
)

// init
func init() {

	// create a console logger
	Console = lumber.NewConsoleLogger(lumber.INFO)

	// create a file logger
	if Log, err = lumber.NewTruncateLogger(Root + "/nanobox.log"); err != nil {
		Fatal("[config/log] lumber.NewAppendLogger() failed", err.Error())
	}
}

// Debug
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}

// Info
func Info(msg string, debug bool) {
	Log.Info(msg)
}

// Fatal
func Fatal(msg, err string) {
	fmt.Println("A fatal error occurred (See ~/.nanobox/nanobox.log for details). Exiting...")
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))
	Log.Close()
	os.Exit(1)
}
