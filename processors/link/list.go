package link

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// List ...
func List(envModel *models.Env) error {

	fmt.Printf("%+v\n", envModel.Links)

	return nil
}
