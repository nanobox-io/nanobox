// Package main ...
package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/bugsnag/bugsnag-go"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	// "github.com/nanobox-io/nanobox/util/memory_logger"
)

var bugsnagToken string

type bugLog struct {
}

func (bugLog) Printf(fmt string, v ...interface{}) {
	lumber.Info(fmt, v...)
}

// main
func main() {
	// build the viper config because viper cannot handle concurrency
	// so it has to be done at the beginning even if we dont need it
	config.Viper()

	// setup a file logger, this will be replaced in verbose mode.
	fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
	if err != nil {
		fmt.Println("logging error:", err)
	}

	//
	lumber.SetLogger(fileLogger)
	lumber.Level(lumber.INFO)
	defer lumber.Close()

	// global panic handler; this is done to avoid showing any panic output if
	// something happens to fail. The output is logged and "pretty" message is
	// shown
	defer func() {
		if r := recover(); r != nil {
			// put r into your log ( it contains the panic message)
			// Then log debug.Stack (from the runtime/debug package)

			stack := debug.Stack()

			bugsnag.Notify(fmt.Errorf("panic"), bugsnag.SeverityError, bugsnag.User{Id: util.UniqueID()}, r, stack)

			lumber.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			lumber.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			lumber.Close()
			fmt.Println("Nanobox encountered an unexpected error. Please see ~/.nanobox/nanobox.log and submit the issue to us.")
			os.Exit(1)
		}
	}()

	setupBugsnag()

	//
	commands.NanoboxCmd.Execute()
}

func setupBugsnag() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:      bugsnagToken,
		Logger:      bugLog{},
		Synchronous: true,
	})

	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		// set the grouping hash to a md5 of the message. which should seperate the grouping in the dashboard
		event.GroupingHash = fmt.Sprintf("%x", md5.Sum([]byte(event.Message)))
		return nil
	})
}
