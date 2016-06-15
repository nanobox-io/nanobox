package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util/data"
)

type processLogout struct {
	control ProcessControl
}

//
func init() {
	Register("logout", logoutFn)
}

//
func logoutFn(conf ProcessControl) (Processor, error) {
	return processLogout{conf}, nil
}

//
func (logout processLogout) Results() ProcessControl {
	return logout.control
}

//
func (logout processLogout) Process() error {
	// remove token from database
	err := data.Delete("global", "user")
	if err == nil {
		fmt.Println("logout successful.")
	}

	return err
}
