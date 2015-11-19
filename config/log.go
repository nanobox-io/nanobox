//
package config

import (
	"fmt"
	"github.com/jcelliott/lumber"
	"os"
)

var (
	Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
	LogFile string
)

// init
func init() {

	// set log file
	LogFile = Root + "/nanobox.log"

	// create a console logger
	Console = lumber.NewConsoleLogger(lumber.INFO)

	// create a file logger
	if Log, err = lumber.NewAppendLogger(LogFile); err != nil {
		Fatal("[config/log] lumber.NewAppendLogger() failed", err.Error())
	}
}

// Debug
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}

// Error
func Error(msg, err string) {
	fmt.Printf("%s (See %s for details)\n", msg, LogFile)
	Log.Error(err)
}

// Fatal
func Fatal(msg, err string) {
	fmt.Printf("A Vagrant error occurred (See %s for details). Exiting...", LogFile)
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))
	Log.Close()
	os.Exit(1)
}
