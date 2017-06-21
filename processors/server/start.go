package server

import (
	"time"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/service"
)

func Start() error {
	// run as admin
	// the reExecPrivilageStart function is defined in the setup
	// since the service create is idempotent it is fine to only have one
	// start command for the server
	if !util.IsPrivileged() {
		return reExecPrivilageStart()
	}

	fn := func() error {
		return service.Start("nanobox-server")
	}

	return util.Retry(fn, 3, time.Second)
}
