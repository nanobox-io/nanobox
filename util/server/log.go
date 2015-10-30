// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/config"
)

//
var (
	Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
)

// NewLogger
func NewLogger(path string) {

	var err error

	// create a console logger
	Console = lumber.NewConsoleLogger(lumber.INFO)

	// create a file logger
	if Log, err = lumber.NewTruncateLogger(path); err != nil {
		config.Error("Failed to create a Server logger", err.Error())
	}
}

// Info
func Info(msg string, debug bool) {
	Log.Info(msg)
}

// Debug
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}

// Fatal
func Fatal(msg, err string) {
	fmt.Printf("Nanobox server errored (See %s for details). Exiting...", config.AppDir+"/server.log")
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))
	Log.Close()
	os.Exit(1)
}
