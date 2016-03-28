package validate

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/util"
)

func init() {
	Register("virtualbox", vboxFunc)
}

func vboxFunc() error {
	// ensure virtualbox exists
	if exists := util.VboxExists(); !exists {
		return fmt.Errorf("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
	}
	return nil
}