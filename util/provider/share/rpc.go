package share

import (
	"github.com/nanobox-io/nanobox/commands/server"
)

type ShareRPC struct{}

type Response struct {
	Message string
	Success bool
}

func init() {
	server.Register(&ShareRPC{})
}
