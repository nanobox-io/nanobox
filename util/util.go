//
package util

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/nanobox-io/nanobox/config"
)

// VboxExists ensure virtualbox is installed; if ever there is a virtualbox package
// this can be moved there
func VboxExists() (exists bool) {
	if config.OS == "windows" {
		exists = os.Getenv("VBOX_MSI_INSTALL_PATH") != ""
	} else if err := exec.Command("which", "vboxmanage").Run(); err == nil {
		exists = true
	}

	return
}

// MD5sMatch determines if a local MD5 matches a remote MD5
func MD5sMatch(localFile, remotePath string) (bool, error) {

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
