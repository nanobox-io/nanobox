// Package main ...
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/update"
	"github.com/nanobox-io/nanobox/models"
)

// main ...
func main() {

	path := ""
	var err error
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		// get the location of the current nanobox
		path, err = exec.LookPath(update.Name)
		if err != nil {
			fmt.Printf("Cannot find %s: %s\n", update.Name, err)
		}
	}

	if !util.IsPrivileged() {

		if runtime.GOOS == "windows" {
			// re-run this command as the administrative user
			fmt.Println()
			fmt.Println("The update process requires Administrator privileges.")
			fmt.Println("Another window will be opened as the Administrator to continue this process.")

			// block here until the user hits enter. It's not ideal, but we need to make
			// sure they see the new window open.
			fmt.Println("Enter to continue:")
			var input string
			fmt.Scanln(&input)

		}

		// make sure the .nanobox folder is created by our user
		models.LoadUpdate()

		cmd := fmt.Sprintf("%s \"%s\"", os.Args[0], path)
		if err := util.PrivilegeExec(cmd); err != nil {
			os.Exit(1)
		}

		// we're done
		return
	}

	// run the update
	err = update.Run(path)
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
