//
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
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

	update, err := updatable()
	if err != nil {
		Config.Error("Unable to determing if updates are available", err.Error())
		return
	}

	// if the md5s don't match or it's been forced, update
	switch {
	case update, config.Force:
		if err := runUpdate(); err != nil {
			if _, ok := err.(*os.LinkError); ok {
				fmt.Println(`Nanobox was unable to update, try again with admin privilege (ex. "sudo nanobox update")`)
			} else {
				Config.Fatal("[commands/update] runUpdate() failed", err.Error())
			}
		}
	default:
		fmt.Printf(stylish.SubBullet("[√] Nanobox is up-to-date"))
	}
}

// Update
func Update() error {

	update, err := updatable()
	if err != nil {
		return fmt.Errorf("Nanobox was unable to determine if updates are available - %s", err.Error())
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

		// don't update by default, assuming they'll just do it manually, prompting
		// again after 14 days
		default:
			fmt.Println("You can manually update at any time with 'nanobox update'.")
			return touchUpdate()

		// if yes continue to update
		case "Yes", "yes", "Y", "y":
			if err := runUpdate(); err != nil {
				if _, ok := err.(*os.LinkError); ok {
					fmt.Println(`Nanobox was unable to update, try again with admin privilege (ex. "sudo nanobox update")`)
				} else {
					return fmt.Errorf("Nanobox was unable to update - %s", err.Error())
				}
			}
		}
	}

	return nil
}

// updatable
func updatable() (bool, error) {

	//
	path, err := osext.Executable()
	if err != nil {
		return false, err
	}

	// check the md5 of the current executing cli against the remote md5;
	// os.Args[0] is used as the final interpolation to determine standard/dev versions
	match, err := Util.MD5sMatch(path, fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%s/%s/%s.md5", config.OS, config.ARCH, filepath.Base(os.Args[0])))
	if err != nil {
		return false, err
	}

	// if the MD5's DONT match we want to update
	return !match, nil
}

// runUpdate attemtps to update using the updater; if it's not available nanobox
// will download it and then run it.
func runUpdate() error {

	//
	path, err := osext.Executable()
	if err != nil {
		return err
	}

	// get the directory of the current executing cli
	dir := filepath.Dir(path)

	// see if the updater is available on PATH
	if _, err := exec.LookPath("nanobox-update"); err != nil {

		tmpFile := filepath.Join(config.TmpDir, "nanobox-update")

		// create a tmp updater in tmp dir
		f, err := os.Create(tmpFile)
		if err != nil {
			return err
		}
		defer f.Close()

		// the updateder is not available and needs to be downloaded
		dl := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/updaters/%s/%s/nanobox-update", config.OS, config.ARCH)

		fmt.Printf("Updater not found. Downloading from %s\n", dl)

		fileutil.Progress(dl, f)

		// ensure updater download matches the remote md5; if the download fails for any
		// reason this md5 should NOT match.
		md5 := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/updaters/%s/%s/nanobox-udpate.md5", config.OS, config.ARCH)
		if _, err = util.MD5sMatch(tmpFile, md5); err != nil {
			return err
		}

		// make new updater executable
		if err := os.Chmod(tmpFile, 0755); err != nil {
			return err
		}

		// move updater to the same location as the cli
		if err = os.Rename(tmpFile, filepath.Join(dir, "nanobox-update")); err != nil {
			return err
		}
	}

	cmd := exec.Command(filepath.Join(dir, "nanobox-update"), "-o", filepath.Base(path))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the updater
	if err := cmd.Run(); err != nil {
		Config.Fatal("[commands/update] exec.Command().Run() failed", err.Error())
	}

	// update the .update file
	return touchUpdate()
}

// touchUpdate updates the mod time on the ~/.nanobox/.update file
func touchUpdate() error {
	return os.Chtimes(config.UpdateFile, time.Now(), time.Now())
}
