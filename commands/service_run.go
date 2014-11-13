package commands

import (
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceRunCommand satisfies the Command interface for running commands on an
// app's service
type ServiceRunCommand struct{}

// Help prints detailed help text for the service run command
func (c *ServiceRunCommand) Help() {
	ui.CPrintln(`
Description:
  Run's a command on an application's service. 'Complex' commands MUST be
  formated inside single quotes ('')

    'ls -la'

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  If [service] is not provided, will run command on first detected code service.

Usage:
  pagoda run [-a app-name] [-s service-name/UID] <COMMAND>
  pagoda service:run [-a app-name] [-s service-name/UID] <COMMAND>

  ex. pagoda run -a app-name -s web1 ls -la

Options:
  -a, --app [app-name]
    The name of the app to which the service belongs

  -s, --service [service-name/UID]
    The name or UID of the service on which the command will run

  -i, --identity-file [path]
    The ssh key file to use
  `)
}

// Run attempts to run a specified command an a designated app service. Takes a
// command flag, which is the command to be run. Takes a service flag to designate
// which serivce to run the command on. If no service is provided, it will select
// the first 'codeable' service it finds and run the command there. Also takes an
// identity file flag to all a user to specify the location of their public SSH
// key if necessary
func (c *ServiceRunCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fCommand string

	// If there's no command, prompt for one
	if len(opts) <= 0 {
		fCommand = ui.Prompt("Please specify a command you wish to run (see help for format): ")

		// We should expect opts[0] to be the service.
	} else {
		fCommand = opts[0]
		opts = opts[1:]
	}

	// parse remaining flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fService string
	flags.StringVar(&fService, "s", "", "")
	flags.StringVar(&fService, "service", "", "")

	var fIdentity string
	flags.StringVar(&fIdentity, "i", "", "")
	flags.StringVar(&fIdentity, "identity-file", "", "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:run", err)
	}

	// if there's no service, loop through all services and get the first
	// 'codeable' service to run the command on.
	if fService == "" {
		services, err := api.GetAppServices(fApp)
		if err != nil {
			fmt.Println("There was a problem getting '%v's' services. See ~/.pagodabox/log.txt for details", fApp)
			ui.Error("pagoda service:run", err)
		}

		for _, service := range services {
			if service.Codeable {
				fService = service.UID
			}
		}
	}

	// the service to run the command on
	service, err := helpers.GetServiceBySlug(fApp, fService, api)
	if err != nil {
		fmt.Printf("Oops! We could not find a '%v' on '%v'.\n", fService, fApp)
		os.Exit(1)
	}

	// the public SSH key to use when SSHing into the server to run the command
	key, err := helpers.GetKeyFile(fIdentity)
	if err != nil {
		fmt.Println("Unable to use SSH key. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:run", err)
	}

	// the config credentials for connecting to the remote server
	config := &ssh.ClientConfig{
		User: service.TunnelUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// all the necessary options for connecting to the remote server
	runOptions := &helpers.SSHOptions{
		Command:    strings.Join([]string{fCommand}, " "),
		Config:     config,
		RemoteIP:   service.TunnelIP,
		RemotePort: service.TunnelPort,
		RemoteUser: service.TunnelUser,
	}

	// if the service SSH tunnel isn't enabled, enable it
	if !service.PublicTunnel {
		helpers.EnablePublicTunnel(service, api, runOptions)
	}

	// create an SSH client
	client, err := ssh.Dial("tcp", runOptions.RemoteIP+":"+strconv.Itoa(runOptions.RemotePort), runOptions.Config)
	if err != nil {
		fmt.Println("Unable to connect SSH client. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:run", err)
	}

	defer client.Close()

	// create an SSH session for client
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Unable to create new SSH client session. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:run", err)
	}

	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 24, 80, modes); err != nil {
		fmt.Println("Unable to request PTY. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:run", err)
	}

	fmt.Println(`
Running '` + runOptions.Command + `' on ` + service.UID + ` :
---------------------------------------------------------
  `)

	if err := session.Run("source /home/gopagoda/.bashrc && cd ./data && " + runOptions.Command); err != nil {
		fmt.Printf("Unable to run command: %v \n", err)
		os.Exit(1)
	}
}
