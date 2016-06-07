package processor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/kardianos/osext"
	cryptoutil "github.com/sdomino/go-util/crypto"
	fileutil "github.com/sdomino/go-util/file"
	printutil "github.com/sdomino/go-util/print"

	"github.com/nanobox-io/nanobox/util"
)

//
var (
	pathToDownload = "https://s3.amazonaws.com/tools.nanobox.io/nanobox/v1"
)

// update is the process created for updating nanobox
type update struct {
	control ProcessControl
}

// init registers "updateFunc" to "update" command
func init() {
	Register("update", updateFunc)
}

// updateFunc returns "update" (above) to be run as a processor
func updateFunc(control ProcessControl) (Processor, error) {
	return &update{control}, nil
}

// Results ...
func (u update) Results() ProcessControl {
	return u.control
}

// Process ...
func (u *update) Process() error {

	// determine if nanobox needs to be updated
	update, err := updatable()
	if err != nil {
		return fmt.Errorf("Unable to determin if updates are available %v", err.Error())
	}

	// NOTE: we just want to os.Exit(0) after an update; the reason is contextual...
	// If a user was trying to run a command and got thrown into an update process
	// to then be thrown back into their command could be confusing, so we'll just
	// exit and let them re-run their command
	switch {

	// update if the update command is run and updates are available
	case update:
		manualUpdate()
		os.Exit(0)

	// update if the update command is run, updates are available, AND the command
	// as forced; this is how we handle our auto update...
	case update && u.control.Force:
		autoUpdate()
		os.Exit(0)
	}

	return nil
}

// updatable determines if the local version of nanobox matches the published version;
// if they don't match then nanobox needs to update
func updatable() (bool, error) {

	// get the path the current executing nanobox
	path, err := osext.Executable()
	if err != nil {
		return false, err
	}

	// check the md5 of the current nanobox against the remote md5; os.Args[0] is
	// used as the final interpolation to determine standard/dev versions
	match, err := cryptoutil.MD5Match(path, fmt.Sprintf("%s/%s/%s/%s.md5", pathToDownload, runtime.GOOS, runtime.GOARCH, filepath.Base(os.Args[0])))
	if err != nil {
		return false, err
	}

	// if the MD5's DONT match we want to update
	return !match, nil
}

// autoUpdate prompts the user if they would like to update nanobox
func autoUpdate() error {

	//
	switch printutil.Prompt("Nanobox is out of date, would you like to update it now (y/N)? ") {

	// don't update by default, assuming they'll just do it manually, prompting
	// again after 14 days
	default:
		fmt.Println("You can manually update at any time with 'nanobox update'.")
		return touchUpdate()

	// if yes continue to update
	case "Yes", "yes", "Y", "y":
		manualUpdate()
	}

	return nil
}

// manualUpdate attempts to update nanobox
func manualUpdate() {
	if err := runUpdate(); err != nil {
		if _, ok := err.(*os.LinkError); ok {
			fmt.Println(`Nanobox was unable to update, try again with admin privilege (ex. "sudo nanobox update")`)
		} else {
			fmt.Println("Nanobox failed to update - ", err.Error())
		}
	}
}

// runUpdate attemtps to update nanobox using the updater; if it's not available
// nanobox will download it to the same location it is running at and then run it.
func runUpdate() error {

	//
	path, err := osext.Executable()
	if err != nil {
		return err
	}

	// get the directory of the current executing cli
	dir := filepath.Dir(path)

	// see if the updater is available on PATH
	if _, err := exec.LookPath("nanobox-updater"); err != nil {
		if err := downloadUpdater(dir); err != nil {
			return err
		}
	}

	// updater command
	cmd := exec.Command(filepath.Join(dir, "nanobox-updater"), "-o", filepath.Base(path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the updater
	if err := cmd.Run(); err != nil {
		lumber.Fatal("[commands/update] exec.Command().Run() failed", err.Error())
	}

	// update the .update file
	return touchUpdate()
}

// downloadUpdater attempts to download the nanobox-updater from S3.
func downloadUpdater(location string) error {
	tmpFile := filepath.Join(util.TmpDir(), "nanobox-updater")

	// create a tmp updater in tmp dir
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// the updateder is not available and needs to be downloaded
	dl := fmt.Sprintf("%s/%s/%s/nanobox-updater", pathToDownload, runtime.GOOS, runtime.GOARCH)

	fmt.Printf("'nanobox-updater' not found. Downloading from %s\n", dl)

	// download the updater
	fileutil.Progress(dl, f)

	// ensure updater download matches the remote md5; if the download fails for
	// any reason this md5 should NOT match.
	md5 := fmt.Sprintf("%s/%s/%s/nanobox-updater.md5", pathToDownload, runtime.GOOS, runtime.GOARCH)
	if _, err = cryptoutil.MD5Match(tmpFile, md5); err != nil {
		return err
	}

	// make new updater executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		return err
	}

	// move updater to the same location as the cli
	if err = os.Rename(tmpFile, filepath.Join(location, "nanobox-updater")); err != nil {
		return err
	}

	return nil
}

// touchUpdate updates the mod time on the ~/.nanobox/.update file
func touchUpdate() error {
	return os.Chtimes(util.UpdateFile(), time.Now(), time.Now())
}
