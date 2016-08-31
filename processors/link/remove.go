package link

import (
	"github.com/nanobox-io/nanobox/models"
)

// Remove ...
func Remove(envModel *models.Env, alias string) error {

	//
	delete(envModel.Links, alias)

	return envModel.Save()
}
