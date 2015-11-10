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
	if Log, err = lumber.NewAppendLogger(path); err != nil {
		config.Fatal("[util/server/log] lumber.NewAppendLogger() failed", err.Error())
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
