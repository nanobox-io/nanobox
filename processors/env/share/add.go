package share

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/provider/share"
)

// Add adds a share share to the workstation
func Add(path string) error {

	// since we dont 
	// // short-circuit if the entry already exist
	// if share.Exists(path) {
	// 	return nil
	// }

	// add the share entry
	if err := share.Add(path); err != nil {
		lumber.Error("share:Add:share.Add(%s): %s", path, err.Error())
		return util.ErrorAppend(err, "failed to add share")
	}

	return nil
}
