//
package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/server"
)

//
var watchCmd = &cobra.Command{
	Hidden: true,

	Use:   "watch",
	Short: "",
	Long:  ``,

	Run: watch,
}

// watch
func watch(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Watching app files for changes"))

	// begin watching for file changes at cwd
	if err := notify.Watch(config.CWDir, server.NotifyRebuild); err != nil {
		fmt.Printf(err.Error())
	}
}
