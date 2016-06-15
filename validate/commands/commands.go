// Package commands ...
package commands

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/validate"
	"github.com/spf13/cobra"
)

// Validate ...
func Validate(checks ...string) func(ccmd *cobra.Command, args []string) {
	return func(ccmd *cobra.Command, args []string) {
		if err := validate.Check(checks...); err != nil {
			fmt.Printf("Validation Failed:\n%s\n", err.Error())
			os.Exit(1)
		}
	}
}
