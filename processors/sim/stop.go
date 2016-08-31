package sim

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

// Stop ...
func Stop(appModel *models.App) error {
	return app.Stop(appModel)
}
