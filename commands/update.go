//
package commands

import (
	"fmt"
	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	printutil "github.com/nanobox-io/nanobox/util/print"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the CLI to the newest available version",
	Long:  ``,

	Run: update,
}

// update
func update(ccmd *cobra.Command, args []string) {
	Update()
	fmt.Printf("Nanobox is now up to date (running v%s)\n", config.VERSION)
}

// Update
func Update() {

	// if there is no nanobox.md5 this is likely the first time nanobox is being
	// run and it should update
	if fi, _ := os.Stat(config.Root + "/nanobox.md5"); fi == nil {
		fmt.Printf(stylish.Bullet("Ensuring nanobox is up to date..."))
		runUpdate()
		return
	}

	//
	match, err := Util.MD5sMatch(config.Root+"/nanobox.md5", "https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5")
	if err != nil {
		Config.Fatal("[commands/update] util.MD5sMatch() failed", err.Error())
	}

	// an error here just means the file doesn't exist (which should never happen
	// since it gets created in the config init at startup)
	fi, _ := os.Stat(config.UpdateFile)

	// if the last update was longer ago than our wait time (14 days), and the md5s
	// dont match, then update
	if time.Since(fi.ModTime()).Hours() >= 336.0 && !match {

		//
		switch printutil.Prompt("Nanobox is out of date, would you like to update it now (y/N)? ") {

		// don't update by default
		default:
			fmt.Println("You can manually update at any time with 'nanobox update'.")
			return

		// if yes continue to update
		case "Yes", "yes", "Y", "y":
			runUpdate()
		}
	}
}

// runUpdate
func runUpdate() {

	fmt.Printf(stylish.Bullet("Updating nanobox..."))

	// get the path of the current executing CLI
	path, err := osext.Executable()
	if err != nil {
		Config.Fatal("[commands/update] osext.ExecutableFolder() failed", err.Error())
	}

	// download the CLI
	cli, err := os.Create(path)
	if err != nil {
		Config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer cli.Close()

	//
	fileutil.Progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH), cli)

	//
	// download the CLI md5
	md5, err := os.Create(config.Root + "/nanobox.md5")
	if err != nil {
		Config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer md5.Close()

	//
	fileutil.Download("https://s3.amazonaws.com/tools.nanobox.io/cli/nanobox.md5", md5)

	// if the new CLI fails to execute, just print a generic message and return
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		fmt.Printf(stylish.SubBullet("[√] Update successful"))
		return
	}

	fmt.Printf(stylish.SubBullet("[√] Now running %s", string(out)))

	// update the modification time of the .update file
	if err := os.Chtimes(config.UpdateFile, time.Now(), time.Now()); err != nil {
		Config.Fatal("[commands.update] os.Chtimes() failed", err.Error())
	}
}
