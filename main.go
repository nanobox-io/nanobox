// Nanobox automates the creation of isolated, repeatable environments for local
// and production applications. When developing locally, Nanobox provisions your
// app's infrastructure inside of a virtual machine (VM) and mounts your local
// codebase into the VM. Any changes made to your codebase are reflected inside
// the virtual environment.
//
// Once code is built and tested locally, Nanobox provisions and deploys an
// identical infrastructure on a production platform.
package main

import (
	"bufio"
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
	proc_provider "github.com/nanobox-io/nanobox/processors/provider"
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
	// setup a file logger, this will be replaced in verbose mode.
	fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
	if err != nil {
		fmt.Println("logging error:", err)
	}

	//
	lumber.SetLogger(fileLogger)
	lumber.Level(lumber.INFO)
	defer lumber.Close()

	// if it is running the server just run it
	// skip the tratiotional messaging
	if len(os.Args) == 2 && (os.Args[1] == "server" || os.Args[1] == "version" || os.Args[1] == "tunnel" || os.Args[1] == "login" || os.Args[1] == "logout") {
		err := commands.NanoboxCmd.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	// verify that we support the prompt they are using
	if badTerminal() {
		display.BadTerminal()
		os.Exit(1)
	}

	// do the commands configure check here because we need it to happen before setupBugsnag creates the config
	command := strings.Join(os.Args, " ")
	if _, err := models.LoadConfig(); err != nil && !strings.Contains(command, " config") && !strings.Contains(command, "env server") {
		processors.Configure()
	}

	migrationCheck()

	fixRunArgs()

	configModel, _ := models.LoadConfig()

	// build the viper config because viper cannot handle concurrency
	// so it has to be done at the beginning even if we dont need it
	providerName := configModel.Provider

	// make sure nanobox has all the necessry parts
	if !strings.Contains(command, " config") && !strings.Contains(command, " server") {
		valid, missingParts := provider.Valid()
		if !valid {

			display.MissingDependencies(providerName, missingParts)
			os.Exit(1)
		}
	}

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
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       bugsnagToken,
		Logger:       bugLog{},
		Synchronous:  true,
		AppVersion:   models.VersionString(),
		PanicHandler: func() {}, // the built in panic handler reexicutes our code
	})

	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		// set the grouping hash to a md5 of the message. which should seperate the grouping in the dashboard
		event.GroupingHash = fmt.Sprintf("%v", event.Stacktrace)
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

	// unset all the docer variables and re init the docker client
	os.Unsetenv("DOCKER_TLS_VERIFY")
	os.Unsetenv("DOCKER_MACHINE_NAME")
	os.Unsetenv("DOCKER_HOST")
	os.Unsetenv("DOCKER_CERT_PATH")

	if err := proc_provider.Init(); err != nil {
		os.Exit(0)
	}

}
