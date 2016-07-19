package validate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Requires ...
func Requires(checks ...string) func(ccmd *cobra.Command, args []string) {

	//
	return func(ccmd *cobra.Command, args []string) {

		//
		if err := Check(checks...); err != nil {
			fmt.Printf("Missing dependencies:\n%s", err.Error())
			os.Exit(1)
		}
	}
}
