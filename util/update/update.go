// Package update handles the updating of nanobox cli
package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

func remotePath() string {
	return fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/nanobox/v2/%s/%s/%s", runtime.GOOS, runtime.GOARCH, Name)
}

func latestVersion() string {
	remotePath := "https://s3.amazonaws.com/tools.nanobox.io/nanobox/v2/version"
	res, err := http.Get(remotePath)
	if err != nil {
		lumber.Error("update:http.Get(%s): %s", remotePath, err)
		return ""
	}

	// read the remote version string
	vers, err := ioutil.ReadAll(res.Body)
	if err != nil {
		lumber.Error("update:ioutil.ReadAll(body): %s", err)
		return ""
	}
	defer res.Body.Close()

	return string(vers)
}

func newUpdate() models.Update {
	return models.Update{
		LastCheckAt:   time.Now(),
		LastUpdatedAt: time.Now(),
	}
}
