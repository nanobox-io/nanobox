// Package update ...
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

func RemotePath() string {
	return fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/nanobox/v1/%s/%s/%s", runtime.GOOS, runtime.GOARCH, name)
}

func RemoveMd5() string {
	remotePath := RemotePath() + ".md5"

	res, err := http.Get(remotePath)
	if err != nil {
		lumber.Error("update:http.Get(%s): %s", remotePath, err)
		return ""
	}

	// read the remote md5 checksum
	md5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		lumber.Error("update:ioutil.ReadAll(body): %s", err)
		return ""
	}
	defer res.Body.Close()

	return string(md5)
}

func populateUpdate(update *models.Update) {
	update.CurrentVersion = RemoveMd5()
	update.LastCheckAt = time.Now()
	update.LastUpdatedAt = time.Now()
}
