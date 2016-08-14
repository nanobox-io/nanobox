package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
type List struct {
	Env models.Env
}

//
func (list List) Run() error {

	// store the auth token
	fmt.Printf("%+v\n", list.Env.Links)

	return nil
}
