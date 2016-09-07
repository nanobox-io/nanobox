// Package main ...
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/jcelliott/lumber"
	"github.com/timehop/go-mixpanel"

	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/update"
)

var bugsnagToken string
var mixpanelToken string

type mixLog struct {
}

func (mixLog) Printf(fmt string, v ...interface{}) {
	lumber.Info(fmt, v...)
}

// main
func main() {

	// setup a file logger, this will be replaced in verbose mode.
	fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
	if err != nil {
		fmt.Println("logging error:", err)
	}

	//
	lumber.SetLogger(fileLogger)
	lumber.Level(lumber.INFO)
	defer lumber.Close()

	setupBugsnag()
	go reportMixpanel()

	// global panic handler; this is done to avoid showing any panic output if
	// something happens to fail. The output is logged and "pretty" message is
	// shown
	defer func() {
		if r := recover(); r != nil {
			// put r into your log ( it contains the panic message)
			// Then log debug.Stack (from the runtime/debug package)

			stack := debug.Stack()

			bugsnag.Notify(fmt.Errorf("panic"), bugsnag.SeverityError, bugsnag.User{Id: config.Viper().GetString("token")}, r, stack)

			lumber.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			lumber.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			lumber.Close()
			fmt.Println("Nanobox encountered an unexpected error. Please see ~/.nanobox/nanobox.log and submit the issue to us.")
			os.Exit(1)
		}
	}()

	// check to see if nanobox needs to be updated
	if err := update.Check(); err != nil {
		fmt.Println("Nanobox was unable to update because of the following error:\n", err.Error())
	}

	//
	commands.NanoboxCmd.Execute()
}

func setupBugsnag() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:      bugsnagToken,
		Logger:      mixLog{},
		Synchronous: true,
	})
}

func reportMixpanel() {
	mx := mixpanel.NewMixpanel(mixpanelToken)

	args := strings.Join(os.Args[1:], " ")
	token := config.Viper().GetString("token")

	err := mx.Track(token, "command", mixpanel.Properties{
		"os":         runtime.GOOS,
		"provider":   config.Viper().GetString("provider"),
		"mount-type": config.Viper().GetString("mount-type"),
		"args":       args,
		"cpus":       runtime.NumCPU(),
	})
	if err != nil {
		lumber.Error("reportMixpanel(): %s", err)
	}
}
