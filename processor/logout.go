package processor

import (
	"fmt"
	"github.com/nanobox-io/nanobox/util/data"
)

type logout struct {
	control ProcessControl
}

func init() {
	Register("logout", logoutFunc)
}

func logoutFunc(conf ProcessControl) (Processor, error) {
	return logout{conf}, nil
}

func (self logout) Results() ProcessControl {
	return self.control
}

func (self logout) Process() error {
	// remove token from database
	err := data.Delete("global", "user")
	if err == nil {
		fmt.Println("logout successful.")
	}

	return err
}
