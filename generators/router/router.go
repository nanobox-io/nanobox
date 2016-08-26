package router

import (

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox-boxfile"
)

//
func loadBoxfile(appModel *models.App) boxfile.Boxfile {
	return boxfile.New([]byte(appModel.DeployedBoxfile))
}
