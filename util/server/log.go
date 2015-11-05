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
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/config"
	"os"
)

//
var (
	Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
)

// create a console and default file logger
func init() {

	// create a default console logger
	Console = config.Console

	// create a default file logger
	Log = config.Log
}

// NewLogger sets the vagrant logger to the given path
func NewLogger(path string) {

	var err error

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
