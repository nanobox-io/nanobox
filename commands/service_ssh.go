package commands

import (
	"bufio"
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"fmt"
	"os"
	"strconv"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceSSHCommand satisfies the Command interface for SSHing into an app's
// service
type ServiceSSHCommand struct{}

// Help prints detailed help text for the service ssh command
func (c *ServiceSSHCommand) Help() {
	ui.CPrintln(`
Description:
  SSH into an application's service

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  The service-name/UID is [yellow]required[reset]. (Ex. web1)

Usage:
  pagoda ssh [-a app-name] service-name/UID
  pagoda service:ssh [-a app-name] service-name/UID

  pagoda ssh -a app-name web1

Options:
  -a, --app [app-name]
    The name of the app to which the service belongs

  -i, --identity-file [path]
    The ssh key file to use
  `)
}

// Run attempts to SSH into an app service and provide a pseudo terminal to use
// on the server. Also takes an identity file flag allowing a user to specify the location
// of their public SSH key if necessary
func (c *ServiceSSHCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fService string

	// If there's no service, prompt for one
	if len(opts) <= 0 {
		fService = ui.Prompt("Into which service would you like to SSH: ")

		// We should expect opts[0] to be the service.
	} else {
		fService = opts[0]
		opts = opts[1:]
	}

	// parse remaining flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fIdentity string
	flags.StringVar(&fIdentity, "i", "", "")
	flags.StringVar(&fIdentity, "identity-file", "", "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:tunnel", err)
	}

	// the service to SSH into
	service, err := helpers.GetServiceBySlug(fApp, fService, api)
	if err != nil {
		fmt.Printf("Oops! We could not find a '%v' on '%v'.\n", fService, fApp)
		os.Exit(1)
	}

	// the public SSH key to use when SSHing into the server
	key, err := helpers.GetKeyFile(fIdentity)
	if err != nil {
		fmt.Println("Unable to use SSH key. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:ssh", err)
	}

	// the config credentials for connecting to the remote server
	config := &ssh.ClientConfig{
		User: service.TunnelUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// all the necessary options for connecting to the remote server
	sshOptions := &helpers.SSHOptions{
		Config:     config,
		RemoteIP:   service.TunnelIP,
		RemotePort: service.TunnelPort,
		RemoteUser: service.TunnelUser,
	}

	// if the service SSH tunnel isn't enabled, enable it
	if !service.PublicTunnel {
		helpers.EnablePublicTunnel(service, api, sshOptions)
	}

	// create an SSH client
	client, err := ssh.Dial("tcp", sshOptions.RemoteIP+":"+strconv.Itoa(sshOptions.RemotePort), sshOptions.Config)
	if err != nil {
		fmt.Println("Unable to connect SSH client. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:ssh", err)
	}

	defer client.Close()

	// create an SSH session for the client
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Unable to create new SSH client session. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda", err)
	}

	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println("Unable to create session pipe. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:ssh", err)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 24, 80, modes); err != nil {
		fmt.Printf("Request for PTY failed: %v \n", err)
		os.Exit(1)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		fmt.Printf("Failed to start shell: %v", err)
		os.Exit(1)
	}

	fmt.Println(`
SSH session established, use ctrl-c to terminate session.
---------------------------------------------------------
  `)

	// Accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda service:ssh", err)
		}

		fmt.Fprint(stdin, str)
	}
}
