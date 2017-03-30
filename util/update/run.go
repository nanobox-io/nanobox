package update

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"io/ioutil"

	"github.com/nanobox-io/nanobox/models"
	// "github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var	Server bool

func Run(path string) error {
	if path == "" {
		fmt.Errorf("invalid path")
	}

	// create a temporary file
	tmpFileName := filepath.Join(filepath.Dir(path), TmpName)
	tmpFile, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	// download the file and display the progress bar
	resp, err := http.Get(RemotePath())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dp := display.DownloadPercent{Total: resp.ContentLength}
	if Server {
		// on the Server we dont really care to see this
		dp.Output = ioutil.Discard
	}
	dp.Copy(tmpFile, resp.Body)

	// close the tmp file
	tmpFile.Close()

	// replace binary
	if err := os.Rename(tmpFileName, path); err != nil {
		return err
	}

	// update the model
	update, _ := models.LoadUpdate()
	populateUpdate(update)

	return update.Save()
}
