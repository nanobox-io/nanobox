//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

//
var resumeCmd = &cobra.Command{
	Hidden: true,

	Use:   "resume",
	Short: "Resumes the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    resume,
}

// resume runs 'vagrant resume'
func resume(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	fmt.Printf(stylish.Bullet("Resuming nanobox..."))
	if err := Vagrant.Resume(); err != nil {
		Config.Fatal("[commands/resume] vagrant.Resume() failed", err.Error())
	}
}
