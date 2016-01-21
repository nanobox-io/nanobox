//
package commands

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/server"
)

var updateImagesCmd = &cobra.Command{
	Use:   "update-images",
	Short: "Updates the nanobox docker images",
	Long:  ``,

	PreRun:  boot,
	Run:     updateImages,
	PostRun: halt,
}

// updateImages
func updateImages(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	// stream update output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "imageupdate"}, Mist.ImageUpdates)
	}()

	fmt.Printf(stylish.Bullet("Updating nanobox docker images (this may take a while)..."))

	// run an image update
	if err := server.Update(""); err != nil {
		server.Fatal("[commands/update-images] server.Update() failed", err.Error())
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
