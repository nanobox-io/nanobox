package app

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dns"
)

// Info ...
func Info(env *models.Env, app *models.App) error {

	if app.State != "active" {
		fmt.Printf("\n------------------------------------------------------\n")
		fmt.Printf("Whoops, it doesn't look like a dev environment has\n")
		fmt.Printf("been setup for the current working directory:\n")
		fmt.Printf("\n%s\n", config.LocalDir())
		fmt.Printf("------------------------------------------------------\n\n")

		return nil
	}

	// print header
	line := strings.Repeat("-", len(env.Name)+32)
	fmt.Printf("\n%s\n", line)
	fmt.Printf("%s (%s)              Status: %s  \n", env.Name, app.Name, app.Status)
	fmt.Printf("%s\n", line)

	fmt.Printf("\nMount Path: %s\n", env.Directory)
	fmt.Printf("Env IP: %s\n", app.GlobalIPs["env"])

	components, _ := app.Components()

	for _, component := range components {

		// print the component header
		if component.Name != component.Label {
			fmt.Printf("\n%s (%s)\n", component.Name, component.Label)
		} else {
			fmt.Printf("\n%s\n", component.Name)
		}

		// print the IP
		fmt.Printf("  IP      : %s\n", component.ExternalIP)

		// print users
		if len(component.Plan.Users) > 0 {
			fmt.Printf("  User(s) :\n")
			for _, user := range component.Plan.Users {
				fmt.Printf("    %s - %s\n", user.Username, user.Password)
			}
		}
	}

	// print environment variables
	fmt.Printf("\nEnvironment Variables\n")
	for key, val := range app.Evars {
		fmt.Printf("  %s = %s\n", key, val)
	}

	// print aliases
	fmt.Printf("\nDNS Aliases\n")
	entries := dns.List(app.ID)

	if len(entries) == 0 {
		fmt.Printf("  none\n")
	} else {
		for _, entry := range entries {
			fmt.Printf("  %s\n", entry.Domain)
		}
	}

	// end on an empty line
	fmt.Println()

	return nil
}
