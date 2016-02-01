//
package dev

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/server"
)

//
var execCmd = &cobra.Command{
	Hidden: true,

	Use:   "exec",
	Short: "Runs a command from inside your app on the nanobox",
	Long:  ``,

	PreRun:  boot,
	Run:     execute,
	PostRun: halt,
}

// execute
func execute(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	//
	if len(args) == 0 {
		args = append(args, print.Prompt("Please specify a command you wish to exec: "))
	}

	//
	v := url.Values{}

	// if a container is found that matches args[0] then set that as a qparam, and
	// remove it from the argument list
	if isContainer(args) {
		v.Add("container", args[0])
		args = args[1:]
	}
	v.Add("cmd", strings.Join(args, " "))

	//
	fmt.Printf(stylish.Bullet("Executing command in nanobox..."))
	if err := server.Exec(v.Encode()); err != nil {
		config.Error("[commands/exec] server.Exec failed", err.Error())
	}

	// PostRun: halt
}

// isContainer
func isContainer(args []string) bool {

	// fetch services to see if the command is trying to run on a specific container
	var services []server.Service
	if _, err := server.Get("/services", &services); err != nil {
		config.Fatal("[commands/exec] server.Get() failed", err.Error())
	}

	// make an exception for build1, as it wont show up on the list, but will always exist
	if args[0] == "build1" {
		return true
	}

	// look for a match
	for _, service := range services {
		if args[0] == service.Name {
			return true
		}
	}
	return false
}
