//
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

//
var publishCmd = &cobra.Command{
	Hidden: true,

	Use:   "publish",
	Short: "(coming soon)",
	Long:  ``,

	Run: publish,
}

// publish
func publish(ccmd *cobra.Command, args []string) {
	fmt.Println("coming soon (http://production.nanobox.io)")
}
