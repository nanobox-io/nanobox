package commands

import (
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	nanoAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/helpers"
	"github.com/nanobox-core/cli/ui"
)

// ServiceTunnelCommand satisfies the Command interface for opening a port forward
// 'tunnel' to an app's service
type ServiceTunnelCommand struct{}

// Help prints detailed help text for the service tunnel command
func (c *ServiceTunnelCommand) Help() {
	ui.CPrintln(`
Description:
  Create an SSH port forward to the specified service

  If [app-name] is not provided, will attempt to detect [app-name] from git
  remotes. If no app or multiple apps detected, will prompt for [app-name].

  The service-name/UID is [yellow]required[reest]. (Ex. database1)

  If [port] is not provided, will use the default [port] of the service type.

    Examples:
    mysql     : 3306
    mongodb   : 27017
    redis     : 6379
    memcached : 11211

Usage:
  pagoda tunnel [-a app-name] service-name/UID [-p port]
  pagoda service:tunnel [-a app-name] service-name/UID [-p port]

  ex. pagoda tunnel -a app-name web1 -p 3000

Options:
  -a, --app [app-name]
    The name of the app to which the service belongs

  -p, --port [port]
    The local port you want to forward

  -i, --identity-file [path]
    The ssh key file to use
  `)
}

// Run attempts to open a port forward 'tunnel' to an app service. Takes a port
// flag to designate which port the user wants to use as their local port forward
// port. Also takes an identity file flag allowing a user to specify the location
// of their public SSH key if necessary. If the tunnel is able to establish correctly
// it will stay open until closed (ctrl + c)
func (c *ServiceTunnelCommand) Run(fApp string, opts []string, api *nanoAPI.Client) {

	// if no app flag was passed, attempt to find one
	if fApp == "" {
		fApp = helpers.FindPagodaApp()
	}

	var fService string

	// If there's no service, prompt for one
	if len(opts) <= 0 {
		fService = ui.Prompt("To which service would you like to open a tunnel: ")

		// We should expect opts[0] to be the service.
	} else {
		fService = opts[0]
		opts = opts[1:]
	}

	// parse remaining flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fPort int
	flags.IntVar(&fPort, "p", 0, "")
	flags.IntVar(&fPort, "port", 0, "")

	var fIdentity string
	flags.StringVar(&fIdentity, "i", "", "")
	flags.StringVar(&fIdentity, "identity-file", "", "")

	if err := flags.Parse(opts); err != nil {
		fmt.Println("Unable to parse flags. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:tunnel", err)
	}

	// the service to open the port forward 'tunnel' to
	service, err := helpers.GetServiceBySlug(fApp, fService, api)
	if err != nil {
		fmt.Printf("Oops! We could not find a '%v' on '%v'.\n", fService, fApp)
		os.Exit(1)
	}

	//
	if fPort == 0 {
		fPort = 0
	}

	// the public SSH key to use when SSHing into the server
	key, err := helpers.GetKeyFile(fIdentity)
	if err != nil {
		fmt.Println("Unable to use SSH key. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:tunnel", err)
	}

	// the config credentials for connecting to the remote server
	config := &ssh.ClientConfig{
		User: service.TunnelUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// all the necessary options for connecting to the remote server
	tunnelOptions := &helpers.SSHOptions{
		Config: config,

		LocalIP:   "localhost",
		LocalPort: fPort,

		RemoteIP:   service.TunnelIP,
		RemotePort: service.TunnelPort,
		RemoteUser: service.TunnelUser,

		ServerIP:   service.IPs["default"],
		ServerPort: 0,

		ServiceApp:  fApp,
		ServiceUser: service.Usernames["default"],
		ServicePass: service.Passwords["default"],
	}

	// if the service SSH tunnel isn't enabled, enable it
	if !service.PublicTunnel {
		helpers.EnablePublicTunnel(service, api, tunnelOptions)
	}

	// create a local listener that listens for traffic on the tunnel from the
	// remote server
	listener, err := net.Listen("tcp", tunnelOptions.LocalIP+":"+strconv.Itoa(tunnelOptions.LocalPort))
	if err != nil {
		fmt.Printf("Unable to tunnel into this type of service: %v \n", err)
		os.Exit(1)
	}

	fmt.Println(`
Tunnel established, use the following credentials to connect:
-------------------------------------------------------------

Host      : ` + tunnelOptions.LocalIP + `
Port      : ` + strconv.Itoa(tunnelOptions.LocalPort) + `
Username  : ` + tunnelOptions.ServiceUser + `
Password  : ` + tunnelOptions.ServicePass + `
Database  : ` + tunnelOptions.RemoteUser + `

-------------------------------------------------------------
(note : ctrl-c To close this tunnel)
  `)

	// attempt to open a tunnel
	for {
		localConnection, err := listener.Accept()

		if err != nil {
			fmt.Println("Unable to maintain tunnel. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda service:tunnel", err)
		}

		// establish the tunnel
		go forward(localConnection, tunnelOptions)
	}
}

// forward attempts to keep open a port forward 'tunnel' to an app service
func forward(localConnection net.Conn, tunnelOptions *helpers.SSHOptions) {

	// Start a client connection to the SSH server
	sshClient, err := ssh.Dial("tcp", tunnelOptions.RemoteIP+":"+strconv.Itoa(tunnelOptions.RemotePort), tunnelOptions.Config)
	if err != nil {
		fmt.Println("Unable to connect SSH client. See ~/.pagodabox/log.txt for details")
		ui.Error("pagoda service:tunnel", err)
	}

	// Initiate a connection to the addr from the remote host
	sshConnection, err := sshClient.Dial("tcp", tunnelOptions.ServerIP+":"+strconv.Itoa(tunnelOptions.ServerPort))

	// Copy localConnection.Reader to sshConnection.Writer
	go func() {
		defer localConnection.Close()

		_, err = io.Copy(localConnection, sshConnection)
		if err != nil {
			fmt.Println("Unable to copy data. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda service:tunnel", err)
		}
	}()

	// Copy sshConnection.Reader to localConnection.Writer
	go func() {
		defer sshConnection.Close()

		_, err = io.Copy(sshConnection, localConnection)
		if err != nil {
			fmt.Println("Unable to copy data. See ~/.pagodabox/log.txt for details")
			ui.Error("pagoda service:tunnel", err)
		}
	}()
}
