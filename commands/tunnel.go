package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (
	// TunnelCmd handles tunneling to components.
	TunnelCmd = &cobra.Command{
		Use:   "tunnel [dry-run|remote-alias] <component.id>",
		Short: "Create a secure tunnel between your local machine & a live component.",
		Long: `
Creates a secure tunnel between your local machine &
a live component. The tunnel allows you to manage
live data using your local client of choice.
`,
		PreRun: steps.Run("login"),
		Run:    tunnelFn,
	}

	// will contain either a listen port or a listen/destination port (chown style `8080:` would be 8080 for both)
	portMap string

	// tunnelPorts contains the ports used in the tunnel
	tunnelPorts = struct {
		listenPort int
		destPort   int
	}{}
)

//
func init() {
	TunnelCmd.Flags().StringVarP(&portMap, "port", "p", "", "Specify a local[:destination] port to listen on and connect to.")
}

// tunnelFn validates the ports and establishes the tunnel
func tunnelFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 2)

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %d). Run the command again with the
name of the component you would like to tunnel into:

ex: nanobox tunnel <component>

`, len(args))

		return
	}

	if portMap != "" {
		ports := strings.Split(portMap, ":")
		if len(ports) > 2 {
			fmt.Printf(`
Please specify a single port pair (source:dest). You specified '%s'.

`, portMap)
			return
		}

		// ports will always be length 1 if they put anything
		port, err := strconv.Atoi(ports[0])
		if err != nil {
			fmt.Printf(`
Please specify a number for a port to listen on. You specified '%s'.'%s'

`, ports[0], err.Error())
			return
		}

		// validate that it's not a priviledged listen port
		if port < 1024 || port > 65535 {
			fmt.Printf(`
Please specify a number between 1024 and 65535 as a port to listen on. You specified '%d'.

`, port)
			return
		}
		tunnelPorts.listenPort = port

		if len(ports) == 2 {
			// heres the chown style 'fill in the blank' magic
			if ports[1] == "" {
				tunnelPorts.destPort = port
			} else {
				port, err := strconv.Atoi(ports[1])
				if err != nil {
					fmt.Printf(`
Please specify a number for a port to tunnel to. You specified '%s'.

`, ports[1])
					return
				}

				// validate that it's a valid port
				if port < 1 || port > 65535 {
					fmt.Printf(`
Please specify a number between 1 and 65535 as a destination port. You specified '%d'.

`, port)
					return
				}
				tunnelPorts.destPort = port

			}
		}
	}

	switch location {
	case "local":
		fmt.Println("tunneling is not required for local development")
		return
	case "production":
		// set the meta arguments to be used in the processor and run the processor
		tunnelConfig := models.TunnelConfig{
			AppName:    name,
			ListenPort: tunnelPorts.listenPort,
			DestPort:   tunnelPorts.destPort,
			Component:  args[0],
		}

		display.CommandErr(processors.Tunnel(env, tunnelConfig))
	}

}
