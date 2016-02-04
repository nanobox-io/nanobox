//
package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	"github.com/nanobox-io/nanobox/util/server"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
)

//
var buildCmd = &cobra.Command{
	Hidden: true,

	Use:   "build",
	Short: "Rebuilds/compiles your app",
	Long:  ``,

	PreRun:  boot,
	Run:     build,
	PostRun: halt,
}

// build
func build(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Building codebase..."))

	// stream build output
	go mistutil.Stream([]string{"log", "deploy"}, mistutil.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- mistutil.Listen([]string{"job", "build"}, mistutil.BuildUpdates)
	}()

	// remount the engine file at ~/.nanobox/apps/<app>/<engine> so any new scripts
	// are used during the build
	if err := engineutil.RemountLocal(); err != nil {
		config.Debug("No engine mounted (not found locally).")
	}

	// run a build
	if err := server.Build(""); err != nil {
		server.Fatal("[commands/build] server.Build() failed", err.Error())
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
