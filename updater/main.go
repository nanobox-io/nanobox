// Package main ...
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/update"
)

// main ...
func main() {

	if runtime.GOOS == "windows" && !util.IsPrivileged() {
		// re-run this command as the administrative user
		fmt.Println()
		fmt.Println("The update process requires Administrator privileges.")
		fmt.Println("Another window will be opened as the Administrator to continue this process.")

		// block here until the user hits enter. It's not ideal, but we need to make
		// sure they see the new window open.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)

		cmd := fmt.Sprintf("%s", os.Args[0])
		if err := util.PrivilegeExec(cmd); err != nil {
			os.Exit(1)
		}

		// we're done
		return
	}

	// run the update
	err := update.Run()
	if err != nil {
		fmt.Println("error: %s", err)
	}

	if runtime.GOOS == "windows" {
		// The update process was spawned in a separate window, which will
		// close as soon as this command is finished. To ensure they see the
		// message, we need to hold open the process until they hit enter.
		fmt.Println()
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	}
}
