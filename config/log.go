//
package config

import (
	"fmt"
	"runtime"
	"time"

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
func Debug(msg string) {
	if Verbose {
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
	fmt.Printf("A fatal error occurred (See %s for details). Exiting...", LogFile)
	Log.Fatal(fmt.Sprintf("%s - %s", msg, err))

	// add a mutex lock in so that if multiple errors are happening at the same
	// time we dont try closing the log twice
	mutex.Lock()

	Log.Close()
	if runtime.GOOS == "windows" {
		// temporarily deadlock
		<-time.After(time.Hour)
	}
	os.Exit(1)
}
