//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	"github.com/spf13/cobra"
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
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "build"}, Mist.BuildUpdates)
	}()

	// remount the engine file at ~/.nanobox/apps/<app>/<engine> so any new scripts
	// are used during the build
	if err := engineutil.RemountLocal(); err != nil {
		config.Error("[util/vagrant/init] engineutil.RemountLocal() failed", err.Error())
	}

	// run a build
	if err := Server.Build(""); err != nil {
		Config.Fatal("[commands/build] server.Build() failed", err.Error())
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
