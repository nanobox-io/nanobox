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
		Short:  "show element from the nanobox database",
		Long:   ``,
		Run:    inspectFunc,
		Hidden: true,
	}
)

// inspectFunc ...
func inspectFunc(ccmd *cobra.Command, args []string) {
	switch {
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
