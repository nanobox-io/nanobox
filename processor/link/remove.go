package link

import (
	"github.com/nanobox-io/nanobox/models"
)

// Remove ...
type Remove struct {
	Env models.Env
	Alias   string
}

//
func (link Remove) Run() error {

	//
	delete(link.Env.Links, link.Alias)

	return link.Env.Save()
}