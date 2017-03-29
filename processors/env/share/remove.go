package share

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/provider/share"
)

// Remove removes a share share from the workstation
func Remove(path string) error {

	// short-circuit if the entry doesn't exist
	if !share.Exists(path) {
		return nil
	}

	// rm the share entry
	if err := share.Remove(path); err != nil {
		lumber.Error("share:Add:share.Remove(%s): %s", path, err.Error())
		return util.ErrorAppend(err, "failed to remove share share")
	}

	return nil
}