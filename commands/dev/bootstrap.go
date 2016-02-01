//
package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/server"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
)

//
var bootstrapCmd = &cobra.Command{
	Hidden: true,

	Use:   "bootstrap",
	Short: "Runs an engine's bootstrap script - downloads code & launches a nanobox",
	Long:  ``,

	PreRun:  boot,
	Run:     bootstrap,
	PostRun: halt,
}

// bootstrap
func bootstrap(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Bootstrapping codebase..."))

	// stream bootstrap output
	go mistutil.Stream([]string{"log", "deploy"}, mistutil.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- mistutil.Listen([]string{"job", "bootstrap"}, mistutil.BootstrapUpdates)
	}()

	// run a bootstrap
	if err := server.Bootstrap(""); err != nil {
		server.Fatal("[commands/bootstrap] server.Bootstrap() failed", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	//
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// PostRun: halt
}
