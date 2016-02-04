//
package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/server"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
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
	go mistutil.Stream([]string{"log", "deploy"}, mistutil.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- mistutil.Listen([]string{"job", "imageupdate"}, mistutil.ImageUpdates)
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
