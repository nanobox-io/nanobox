package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (

	// UpdateCmd ...
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates the Nanobox CLI to the newest available *minor* version.",
		Long: `
use the other one.. not this one.. just left in to help people get to the right place
		`,
		Run: updateFn,
	}
)

// updateFn ...
func updateFn(ccmd *cobra.Command, args []string) {
	fmt.Println("YO DAWG!! you gonna have to use the nanobox-update thingie ma jig")
}
