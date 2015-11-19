//
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	printutil "github.com/nanobox-io/nanobox/util/print"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the CLI to the newest available version",
	Long:  ``,

	Run: update,
}

// update
func update(ccmd *cobra.Command, args []string) {

	update, err := updateAvailable()
	if err != nil {
		fmt.Println("Unable to determing if updates are available (see log for details).")
		Config.Error("[commands/update] updateAvailable() failed", err.Error())
		return
	}

	// if the md5s don't match or it's been forced, update
	switch {
	case update, config.Force:
		runUpdate()
	default:
		fmt.Printf(stylish.SubBullet("[√] Nanobox is up-to-date"))
	}
}

// Update
func Update() {

	update, err := updateAvailable()
	if err != nil {
		fmt.Println("Unable to determing if updates are available (see log for details).")
		Config.Error("[commands/update] updateAvailable() failed", err.Error())
		return
	}

	// stat the update file to get ModTime(); an error here means the file doesn't
	// exist. This is highly unlikely as the file is created if it doesn't exist
	// each time the CLI is run.
	fi, _ := os.Stat(config.UpdateFile)

	// if the md5s don't match and it's 'time' for an update (14 days), OR a force
	// update is issued, update
	if update && time.Since(fi.ModTime()).Hours() >= 336.0 {

		//
		switch printutil.Prompt("Nanobox is out of date, would you like to update it now (y/N)? ") {

		// don't update by default
		default:
			fmt.Println("You can manually update at any time with 'nanobox update'.")

			// if they don't update, assume then that they'll either do it manually or just
			// wait 14 more days
			if err := touchUpdate(); err != nil {
				fmt.Println("Failed to touch update")
				Config.Error("[commands/update] updateAvailable() failed", err.Error())
			}

			return

		// if yes continue to update
		case "Yes", "yes", "Y", "y":
			runUpdate()
		}
	}
}

// updateAvailable
func updateAvailable() (bool, error) {

	// get the path of the current executing CLI
	exe, err := osext.Executable()
	if err != nil {
		return false, err
	}

	// check the current cli md5 against the remote md5; os.Args[0] is used as the
	// final interpolation to determine standard/dev versions
	md5 := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/%v.md5", runtime.GOOS, runtime.GOARCH, filepath.Base(os.Args[0]))

	match, err := Util.MD5sMatch(exe, md5)
	if err != nil {
		return false, err
	}

	return !match, nil
}

// runUpdate
func runUpdate() error {

	fmt.Printf(stylish.Bullet("Updating nanobox..."))

	// get the path of the current executing CLI
	exe, err := osext.Executable()
	if err != nil {
		return err
	}

	//

	prog := filepath.Base(os.Args[0])
	tmpDir := config.Root + "/tmp"
	tmpFile := tmpDir + "/" + prog

	// create a tmp dir to download the new cli to; don't care about the error here
	// because if the tmp dir already exists we'll just use it
	os.Mkdir(tmpDir, 0755)

	// create a tmp cli in tmp dir
	cli, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer cli.Close()

	// download the new cli
	dl := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/%v", runtime.GOOS, runtime.GOARCH, prog)
	fileutil.Progress(dl, cli)

	// make new CLI executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		return err
	}

	// move new executable over current CLI's location
	if err = os.Rename(tmpFile, exe); err != nil {
		return err
	}

	// if the new CLI fails to execute, just print a generic message and return
	out, err := exec.Command(exe, "-v").Output()
	if err != nil {
		fmt.Printf(stylish.SubBullet("[√] Update successful"))
		return nil
	}

	fmt.Printf(stylish.SubBullet("[√] Now running %s", string(out)))

	// update the .update file
	return touchUpdate()
}

// touchUpdate updates the mod time on the ~/.nanobox/.update file
func touchUpdate() error {
	return os.Chtimes(config.UpdateFile, time.Now(), time.Now())
}
