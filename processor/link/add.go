package link

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/odin"
)

// Add
type Add struct {
	Env   models.Env
	App   string
	Alias string
}

//
func (add Add) Run() error {
	if add.Alias == "" {
		add.Alias = "default"
	}

	// get app id
	app, err := odin.App(add.App)
	if err != nil {
		return err
	}

	add.Env.Links[add.Alias] = app.ID

	return add.Env.Save()
}
