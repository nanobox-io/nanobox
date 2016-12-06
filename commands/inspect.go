package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
)

type (
	anything interface {
	}
)

var (
	// InspectCmd ...
	InspectCmd = &cobra.Command{
		Use:    "inspect",
		Short:  "Show element from the nanobox database.",
		Long:   ``,
		Run:    inspectFunc,
		Hidden: true,
	}
)

// inspectFunc ...
func inspectFunc(ccmd *cobra.Command, args []string) {
	switch {
	case len(args) == 1 && args[0] == "ip-tree":
		showIPTree()
	default:
		fmt.Println("I need to know some data starting point")

	case len(args) == 1:
		showData(models.Inspect(args[0], ""))
	case len(args) == 2:
		showData(models.Inspect(args[0], args[1]))
	}
}

func showData(v interface{}) {
	fmt.Printf("%+v\n", v)
}

func showIPTree() {
	fmt.Print("reservedIPS: ")
	showData(models.Inspect("registry", "ips"))
	envs, _ := models.AllEnvs()
	for _, env := range envs {
		fmt.Printf("%s\n", env.Name)
		apps, _ := env.Apps()
		for _, app := range apps {
			fmt.Printf("  %s global: %v local: %v\n", app.Name, app.GlobalIPs, app.LocalIPs)
			components, _ := app.Components()
			for _, component := range components {
				fmt.Printf("    %-15s external: %s internal: %s\n", component.Name, component.ExternalIP, component.InternalIP)
			}
		}
	}
}
