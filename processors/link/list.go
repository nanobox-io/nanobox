package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
func List(envModel *models.Env) error {

	// store the auth token
	fmt.Printf("%+v\n", envModel.Links)

	return nil
}
