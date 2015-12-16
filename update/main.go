//
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// main
func main() {

	var download string

	// accept a flag allowing for alternate download options
	flag.StringVar(&download, "o", "nanobox", "The version of nanobox to update")
	flag.Parse()

	// if download is not one of our available download options default to "nanobox"
	if _, ok := map[string]int{"nanobox": 1, "nanobox-dev": 1}[download]; !ok {
		fmt.Printf("'%s' is not a valid option. Downloading 'nanobox'\n", download)
		download = "nanobox"
	}

	// before attempting to update, ensure nanobox is installed (on the path)
	path, err := exec.LookPath(download)
	if err != nil {
		fmt.Printf("Unable to update '%s' (not found on path)\n", download)
		os.Exit(1)
	}

	fmt.Printf("Updating %s...\n", download)

	// set Home based off the users homedir (~)
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to determine home directory - ", err.Error())
		os.Exit(1)
	}

	tmpDir := filepath.Join(home, ".nanobox", "tmp")
	tmpPath := filepath.Join(tmpDir, download)

	// if tmp dir doesn't exist fail. The updater shouldn't run if nanobox has never
	// been run.
	if _, err := os.Stat(tmpDir); err != nil {
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
	progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/%v", runtime.GOOS, runtime.GOARCH, download), tmpFile)

	// make new CLI executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		fmt.Println("Failed to set permissions - ", err.Error())
	}

	// ensure new CLI download matches the remote md5; if the download fails for any
	// reason these md5's should NOT match.
	if _, err = md5sMatch(tmpPath, fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/%v.md5", runtime.GOOS, runtime.GOARCH, download)); err != nil {
		fmt.Printf("Nanobox was unable to correctly download the update. Please check your internet connection and try again.")
		os.Exit(1)
	}

	// replace the old CLI with the new one
	if err = os.Rename(tmpPath, path); err != nil {
		fmt.Println("Failed to replace existing CLI with new one -", err.Error())
		os.Exit(1)
	}

	// execute the new CLI printing the version to verify update;  if the new CLI
	out, err := exec.Command(path, "-v").Output()

	// fails to execute, just print a generic message and return
	if err != nil {
		fmt.Printf("[√] Update successful")
		return
	}

	fmt.Printf("[√] Now running %s", string(out))
}

// progress downloads a file with a fancy progress bar
func progress(path string, w io.Writer) error {

	//
	download, err := http.Get(path)
	if err != nil {
		return err
	}
	defer download.Body.Close()

	var percent float64
	var down int

	// format the response content length to be more 'friendly'
	total := float64(download.ContentLength) / math.Pow(1024, 2)

	// create a 'buffer' to read into
	p := make([]byte, 2048)

	//
	for {

		// read the response body (streaming)
		n, err := download.Body.Read(p)

		// write to our buffer
		w.Write(p[:n])

		// update the total bytes read
		down += n

		switch {
		default:
			// update the percent downloaded
			percent = (float64(down) / float64(download.ContentLength)) * 100

			// show download progress: 0.0/0.0MB [*** progress *** 0.0%]
			fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

		// if no more files are found return
		case download.ContentLength < 1:
			fmt.Printf("\r   %.2fMB", float64(down)/math.Pow(1024, 2))
		}

		// detect EOF and break the 'stream'
		if err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			} else {
				return err
			}
		}
	}

	return nil
}

// md5sMatch determines if a local MD5 matches a remote MD5
func md5sMatch(localFile, remotePath string) (bool, error) {

	// read the local file; will return os.PathError if doesn't exist
	b, err := ioutil.ReadFile(localFile)
	if err != nil {
		return false, err
	}

	// get local md5 checksum (as a string)
	localMD5 := fmt.Sprintf("%x", md5.Sum(b))

	// GET remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// read the remote md5 checksum
	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// compare checksum's
	return strings.TrimSpace(localMD5) == strings.TrimSpace(string(remoteMD5)), nil
}
