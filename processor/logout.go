package processor

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util/data"
)

type processLogout struct {
	control ProcessControl
}

func init() {
	Register("logout", logoutFn)
}

func logoutFn(conf ProcessControl) (Processor, error) {
	return processLogout{conf}, nil
}

func (logout processLogout) Results() ProcessControl {
	return logout.control
}

// Process ...
func (logout processLogout) Process() error {

	// remove token from database
	if err := logout.deleteUser(); err != nil {
		return err
	}

	fmt.Println("Successfully logged out!")

	return nil
}

// deleteUser ...
func (logout *processLogout) deleteUser() error {

	// remove token from database
	if err := data.Delete("global", "user"); err != nil {
		return err
	}

	return nil
}
