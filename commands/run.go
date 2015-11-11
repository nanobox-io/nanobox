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
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a nanobox, provisions the app, & runs the app's exec",
	Long:  ``,

	PreRun:  boot,
	Run:     run,
	PostRun: halt,
}

//
func init() {
	runCmd.Flags().BoolVarP(&config.Force, "reset-cache", "", false, "resets stuff")
}

// run
func run(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	// stream deploy output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "deploy"}, Mist.DeployUpdates)
	}()

	// remount the engine file at ~/.nanobox/apps/<app>/<engine> so any new scripts
	// will be used during the deploy
	if err := engineutil.RemountLocal(); err != nil {
		config.Error("[util/vagrant/init] engineutil.RemountLocal() failed", err.Error())
	}

	// run a deploy
	if err := Server.Deploy("run=true"); err != nil {
		Config.Fatal("[commands/run] server.Deploy() failed - ", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	//
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf(`
--------------------------------------------------------------------------------
[√] APP SUCCESSFULLY BUILT   ///   DEV URL : %v
--------------------------------------------------------------------------------
`, config.Nanofile.Domain)

	// if in background mode just exist w/o streaming logs or watching files
	if config.VMfile.IsBackground() {
		return
	}

	// if not in background mode begin streaming logs and watching files
	fmt.Printf(`
++> STREAMING LOGS (ctrl-c to exit) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
`)

	// stream app output
	go Mist.Stream([]string{"log", "app"}, Mist.ProcessLogStream)

	// begin watching for file changes (blocking)
	if err := Notify.Watch(config.CWDir, Server.NotifyRebuild); err != nil {
		fmt.Printf(err.Error())
	}

	// PostRun: halt
}
