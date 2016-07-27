// Package main ...
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
	cryptoutil "github.com/sdomino/go-util/crypto"
	fileutil "github.com/sdomino/go-util/file"
	
	"github.com/nanobox-io/nanobox/util"
)

var (
	// path to nanobox download; this is hardcoded since the updater will only be
	// responsible for updating minor and patch versions of nanobox. If another major
	// version is released, it will include it's own downloader
	pathToDownload = "https://s3.amazonaws.com/tools.nanobox.io/nanobox/v1"
	
	// name of the file to download ("nanobox" or "nanobox-dev")
  fileToDownload = "nanobox"

  // map of available downloads
  availableDownloads = map[string]bool{"nanobox": true, "nanobox-dev": true}
)

// main ...
func main() {

	// accept a flag allowing for alternate download options
	flag.StringVar(&fileToDownload, "o", "nanobox", "The version of nanobox to update")
	flag.Parse()

	// if download is not one of our available download options reset to "nanobox"
	if _, ok := availableDownloads[fileToDownload]; !ok {
		fmt.Printf("'%s' is not a valid option. Downloading 'nanobox'\n", fileToDownload)
		fileToDownload = "nanobox"
	}

	// before attempting to update, ensure nanobox is installed (on the path)
	path, err := exec.LookPath(fileToDownload)
	if err != nil {
		fmt.Printf("Unable to update '%s' - %v\n", fileToDownload, err)
		os.Exit(1)
	}

	if runtime.GOOS == "windows" && !util.IsPrivileged() {
		// re-run this command as the administrative user
		fmt.Println("The update process requires Administrator privileges.")
		fmt.Println("Another window will be opened as the Administrator to continue this process.")
		
		// block here until the user hits enter. It's not ideal, but we need to make
		// sure they see the new window open.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
		
		cmd := fmt.Sprintf("%s -o %s", os.Args[0], fileToDownload)
		if err := util.PrivilegeExec(cmd); err != nil {
			os.Exit(1)
		}
		
		// we're done
		return
	}

	// get the current users home dir
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to determine home directory - ", err.Error())
		os.Exit(1)
	}

	// if this is windows, we need to tack an .exe extension onto the file
	if runtime.GOOS == "windows" {
		fileToDownload += ".exe"
	}

	tmpDir := filepath.ToSlash(filepath.Join(home, ".nanobox", "tmp"))
	tmpPath := filepath.ToSlash(filepath.Join(tmpDir, fileToDownload))

	// attempt to make a ~.nanobox/tmp directory just incase it doesn't exist
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fmt.Printf("Failed to create '%v' - %v\n", tmpDir, err.Error())
		os.Exit(1)
	}

	// create a tmp CLI in tmp dir
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		fmt.Println("Failed to create temporary file - ", err.Error())
		os.Exit(1)
	}

	// download the new CLI
	fmt.Printf("Updating %s...\n", fileToDownload)
	pathLabel := filepath.ToSlash(fmt.Sprintf("%s/%s/%s/%s", pathToDownload, runtime.GOOS, runtime.GOARCH, fileToDownload))
	fileutil.Progress(pathLabel, tmpFile)

	// close the handle now so we can move the file later
	tmpFile.Close()

	// ensure new CLI download matches the remote md5; if the download fails for any
	// reason these md5's should NOT match.
	if _, err = cryptoutil.MD5Match(tmpPath, fmt.Sprintf("%s/%s/%s/%s.md5", pathToDownload, runtime.GOOS, runtime.GOARCH, fileToDownload)); err != nil {
		fmt.Printf("Nanobox was unable to correctly download the update. Please check your internet connection and try again.")
		os.Exit(1)
	}

	// The process for windows is different enough than the unixes
	if runtime.GOOS != "windows" {
		// make new CLI executable	
		if err := os.Chmod(tmpPath, 0755); err != nil {
			fmt.Println("Failed to set permissions - ", err.Error())
		}
	}

	// replace the old CLI with the new one
	if err = os.Rename(tmpPath, path); err != nil {
		fmt.Println("Failed to replace existing CLI with new one -", err.Error())
		os.Exit(1)
	}

	// execute the new CLI printing the version to verify update
	out, err := exec.Command(path, "-v").Output()

	// if the new CLI fails to execute, just print a generic message and return
	if err != nil {
		fmt.Printf("[!] Update failed")
	} else {
		fmt.Printf("[√] Now running %s", string(out))
	}

	
	if runtime.GOOS == "windows" {
		// The update process was spawned in a separate window, which will
		// close as soon as this command is finished. To ensure they see the
		// message, we need to hold open the process until they hit enter.
		fmt.Println("Enter to continue:")
		var input string
    fmt.Scanln(&input)
	}
}
