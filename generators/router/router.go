package router

import (
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
)

//
func loadBoxfile(appModel *models.App) boxfile.Boxfile {
	return boxfile.New([]byte(appModel.DeployedBoxfile))
}
