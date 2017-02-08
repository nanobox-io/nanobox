// Package main ...
package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

var bugsnagToken string

type bugLog struct {
}

func (bugLog) Printf(fmt string, v ...interface{}) {
	lumber.Info(fmt, v...)
}

// main
func main() {
	// verify that we support the prompt they are using
	if badTerminal() {
		display.BadTerminal()
		os.Exit(1)
	}

	// do the commands configure check here because we need it to happen before setupBugsnag creates the config
	command := strings.Join(os.Args, " ")
	if _, err := models.LoadConfig(); err != nil && !strings.Contains(command, " config") {
		processors.Configure()
	}

	migrationCheck()

	fixRunArgs()

	configModel, _ := models.LoadConfig()

	// build the viper config because viper cannot handle concurrency
	// so it has to be done at the beginning even if we dont need it
	providerName := configModel.Provider

	// make sure nanobox has all the necessry parts
	if !strings.Contains(command, " config") {
		valid, missingParts := provider.Valid()
		if !valid {
			display.MissingDependencies(providerName, missingParts)
			os.Exit(1)
		}
	}

	// setup a file logger, this will be replaced in verbose mode.
	fileLogger, err := lumber.NewAppendLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
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
			// in a panic state we dont want to try loading or using any non standard libraries.
			// so we will just use the ones we already have
			lumber.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			lumber.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			lumber.Close()
			fmt.Println("Nanobox encountered an unexpected error. Please see ~/.nanobox/nanobox.log and submit the issue to us.")
			os.Exit(1)
		}
	}()

	// get the bugsnag variables ready
	setupBugsnag()

	//
	commands.NanoboxCmd.Execute()
}

func setupBugsnag() {
	update, _ := models.LoadUpdate()
	md5Parts := strings.Fields(update.CurrentVersion)
	version := ""
	if len(md5Parts) > 1 {
		version = md5Parts[len(md5Parts)-1]
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       bugsnagToken,
		Logger:       bugLog{},
		Synchronous:  true,
		AppVersion:   version,
		PanicHandler: func() {}, // the built in panic handler reexicutes our code
	})

	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		// set the grouping hash to a md5 of the message. which should seperate the grouping in the dashboard
		event.GroupingHash = fmt.Sprintf("%x", md5.Sum([]byte(event.Message)))
		return nil
	})

}

func badTerminal() bool {
	return runtime.GOOS == "windows" && strings.Contains(os.Getenv("shell"), "bash")
}

func fixRunArgs() {
	found := false
	lastLocation := 0
LOOP:
	for i, arg := range os.Args {
		switch arg {
		case "run":
			found = true
			lastLocation = i
		case "--debug", "--trace", "--verbose", "-t", "-v":
			// if we hit a argument of ours after 'found'
			// we will reset the last location
			if found == true {
				lastLocation = i
			}
		default:
			// if we hit this after we have found run
			// we are done
			if found == true {
				break LOOP
			}
		}

	}

	if found {
		os.Args = append(os.Args[:lastLocation+1], strings.Join(os.Args[lastLocation+1:], " "))
	}
}

// check to see if we need to wipe the old
func migrationCheck() {
	config, _ := models.LoadConfig()
	providerName := config.Provider
	providerModel, err := models.LoadProvider()

	// if the provider hasnt changed or its a new provider
	// no migration required
	if util.IsPrivileged() || err != nil || providerModel.Name == providerName {
		return
	}

	// remember the new provider
	newProviderName := providerName

	// when migrating from the old system
	// the provider.Name may be blank
	if providerModel.Name == "" {
		display.MigrateOldRequired()
		providerModel.Name = newProviderName
	} else {
		display.MigrateProviderRequired()
	}

	// adjust cached config to be the old provider
	config.Provider = providerModel.Name
	config.Save()

	// alert the user of our actions
	fmt.Println("press enter to continue")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	// implode the old system
	processors.Implode()

	// on implode success
	// adjust the provider to the new one and save the provider model
	config.Provider = newProviderName
	config.Save()

	providerModel.Name = newProviderName
	providerModel.Save()
}
