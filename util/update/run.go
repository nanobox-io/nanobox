package update

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nanobox-io/nanobox/models"
	// "github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

func Run() error {

	// get the location of the current nanobox
	path, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	// create a temporary file

	tmpFileName := filepath.Join(filepath.Dir(path), tmpName)
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
