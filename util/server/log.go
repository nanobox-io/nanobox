//
package server

import (
	"fmt"
	"os"
	"sync"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/config"
)

//
var (
	// Console *lumber.ConsoleLogger
	Log     *lumber.FileLogger
	logFile string
	mutex   = &sync.Mutex{}
)

// create a console and default file logger
func init() {

	// create a default console logger
	// Console = config.Console

	// create a default file logger
	Log = config.Log
	logFile = config.LogFile
}

// NewLogger sets the vagrant logger to the given path
func NewLogger(path string) {

	var err error

	// create a file logger (append if already exists)
	if Log, err = lumber.NewAppendLogger(path); err != nil {
		config.Fatal("[util/server/log] lumber.NewAppendLogger() failed", err.Error())
	}

	logFile = path
}

// Debug
func Debug(msg string) {
	if config.Verbose {
		fmt.Printf(msg)
	}
}

// Error
func Error(msg, err string) {
	fmt.Printf("%s (See %s for details)\n", msg, logFile)
	Log.Error(err)
}

// Fatal
func Fatal(msg, err string) {
	fmt.Printf("A fatal server error occurred (See %s for details). Exiting...\n", logFile)
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))

	// add a mutex lock in so that if multiple errors are happening at the same
	// time we dont try closing the log twice
	mutex.Lock()

	Log.Close()
	os.Exit(1)
}
