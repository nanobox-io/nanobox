// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

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
	logfile string
)

// init
func init() {

	// create a console logger
	Console = lumber.NewConsoleLogger(lumber.INFO)

	// try to use the logfile if the app exists, if not just use the default log
	logfile = config.AppDir + "/vagrant.log"
	if _, err := os.Stat(logfile); err != nil {
		logfile = config.Root + "/nanobox.log"
	}

	// create a file logger
	if Log, err = lumber.NewAppendLogger(logfile); err != nil {
		config.Fatal("[util/vagrant/log] lumber.NewAppendLogger() failed", err.Error())
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
	fmt.Printf("A Vagrant error occurred (See %s for details). Exiting...", logfile)
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))
	Log.Close()
	os.Exit(1)
}
