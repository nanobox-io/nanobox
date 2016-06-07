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
)

var (

	// path to nanobox download; this is hardcoded since the updater will only be
	// responsible for updating minor and patch versions of nanobox. If another major
	// version is released, it will include it's own downloader
	pathToDownload = "https://s3.amazonaws.com/tools.nanobox.io/nanobox/v1"

	// name of the file to download ("nanobox" or "nanobox-dev")
	fileToDownload = "nanobox"

	// map of available downloads
	availableDownloads = map[string]int{"nanobox": 0, "nanobox-dev": 0}
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
		fmt.Printf("Unable to update '%s' (not found on path)\n", fileToDownload)
		os.Exit(1)
	}

	fmt.Printf("Updating %s...\n", fileToDownload)

	// get the current users home dir
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to determine home directory - ", err.Error())
		os.Exit(1)
	}

	tmpDir := filepath.Join(home, ".nanobox", "tmp")
	tmpPath := filepath.Join(tmpDir, fileToDownload)

	// if tmp dir doesn't exist fail. The updater shouldn't run if nanobox has never
	// been run.
	if _, err = os.Stat(tmpDir); err != nil {
		fmt.Println("Nanobox updater required nanobox be configured (run once) before it can update.", err.Error())
		os.Exit(1)
	}

	// create a tmp CLI in tmp dir
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		fmt.Println("Failed to create temporary file - ", err.Error())
		os.Exit(1)
	}
	defer tmpFile.Close()

	// download the new CLI
	fileutil.Progress(fmt.Sprintf("%s/%s/%s/%s", pathToDownload, runtime.GOOS, runtime.GOARCH, fileToDownload), tmpFile)

	// ensure new CLI download matches the remote md5; if the download fails for any
	// reason these md5's should NOT match.
	if _, err = cryptoutil.MD5Match(tmpPath, fmt.Sprintf("%s/%s/%s/%s.md5", pathToDownload, runtime.GOOS, runtime.GOARCH, fileToDownload)); err != nil {
		fmt.Printf("Nanobox was unable to correctly download the update. Please check your internet connection and try again.")
		os.Exit(1)
	}

	// make new CLI executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		fmt.Println("Failed to set permissions - ", err.Error())
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
		fmt.Printf("[√] Update successful")
		return
	}

	fmt.Printf("[√] Now running %s", string(out))
}
