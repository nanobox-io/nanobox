package processor

import (
	"fmt"
	"github.com/nanobox-io/nanobox/util/data"
)

type logout struct {
	config ProcessConfig
}

func init() {
	Register("logout", logoutFunc)
}

func logoutFunc(conf ProcessConfig) (Processor, error) {
	return logout{conf}, nil
}

func (self logout) Results() ProcessConfig {
	return self.config
}

func (self logout) Process() error {
	// remove token from database
	err := data.Delete("global", "user")
	if err == nil {
		fmt.Println("logout successful.")
	}

	return err
}
