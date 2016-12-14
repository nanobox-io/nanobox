// Package main ...
package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"bufio"
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
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/display"

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
		fmt.Println("This console is currently not supported by nanobox")
		fmt.Println("Please refer to the docs for more information")
		os.Exit(1)
	}

	// temp 
	firstVpn()

	fixRunArgs()

	// do the commands configure check here because we need it to happen before setupBugsnag creates the config
	if !config.ConfigExists() {
		processors.Configure()
	}

	// build the viper config because viper cannot handle concurrency
	// so it has to be done at the beginning even if we dont need it
	providerName := config.Viper().GetString("provider")

	// make sure nanobox has all the necessry parts
	valid, missingParts := provider.Valid()
	if !valid {
		fmt.Printf("Using nanobox with %s requires tools that appear to not be available on your system.\n", providerName)
		fmt.Println(strings.Join(missingParts, "\n"))
		if providerName == "native" {
			providerName = "docker"
		}
		fmt.Println("View these requirements at docs.nanobox.io/install/requirements/%s/\n", providerName)
		os.Exit(1)
	}

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
		APIKey:      bugsnagToken,
		Logger:      bugLog{},
		Synchronous: true,
		AppVersion:  version,
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

// a temporary check to help people migrate over to the new vpn way
func firstVpn()  {
	// check to see if this is the first time we are running the new vpn way
	env, _ := models.FindEnvByID(config.EnvID())

	if config.Viper().GetString("provider") == "docker-machine" && (env.ID != "" && env.LastBuildProvider == "") {
		fmt.Println("manditory migration message")
		fmt.Println("press enter to continue")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		display.CommandErr(processors.Implode())

	}

}