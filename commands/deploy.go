//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	"github.com/spf13/cobra"
	"net/url"
	"strconv"
)

var (

	//
	deployCmd = &cobra.Command{
		Hidden: true,

		Use:   "deploy",
		Short: "Deploys code to the nanobox",
		Long:  ``,

		PreRun:  boot,
		Run:     deploy,
		PostRun: halt,
	}

	//
	install bool // tells nanobox server to install services
)

//
func init() {
	deployCmd.Flags().BoolVarP(&install, "run", "", false, "Creates your app environment w/o webs or workers")
}

// deploy
func deploy(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	// stream deploy output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "deploy"}, Mist.DeployUpdates)
	}()

	v := url.Values{}
	v.Add("reset", strconv.FormatBool(config.Force))
	v.Add("run", strconv.FormatBool(install))

	// remount the engine file at ~/.nanobox/apps/<app>/<engine> so any new scripts
	// will be used during the deploy
	if err := engineutil.RemountLocal(); err != nil {
		config.Error("[util/vagrant/init] engineutil.RemountLocal() failed", err.Error())
	}

	// run a deploy
	if err := Server.Deploy(v.Encode()); err != nil {
		Config.Fatal("[commands/deploy] server.Deploy() failed - ", err.Error())
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
