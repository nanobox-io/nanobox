//
package commands

import (
	"github.com/nanobox-io/nanobox/util/server"
	"github.com/spf13/cobra"
	"net/url"
	"strings"
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
		args = append(args, Print.Prompt("Please specify a command you wish to exec: "))
	}

	//
	v := url.Values{}

	// if a container is found that matches args[0] then set that as a qparam, and
	// remove it from the argument list
	if Server.IsContainerExec(args) {
		v.Add("container", args[0])
		args = args[1:]
	}
	v.Add("cmd", strings.Join(args, " "))

	//
	if err := server.Exec("exec", v.Encode()); err != nil {
		server.Error("[commands/exec] Server.Exec failed", err.Error())
	}

	// PostRun: halt
}
