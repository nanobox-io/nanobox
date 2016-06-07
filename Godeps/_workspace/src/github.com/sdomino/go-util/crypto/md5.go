package crypto

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// MD5Match determines if a local MD5 matches a remote MD5
func MD5Match(localFile, remotePath string) (bool, error) {

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
