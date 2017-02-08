package platform

import (
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-portal-client"
	generator "github.com/nanobox-io/nanobox/generators/router"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

// UpdatePortal ...
func UpdatePortal(appModel *models.App) error {
	client := portalClient(appModel)

	// update routes
	routes := generator.BuildRoutes(appModel)
	updateRoute := func() error {
		return client.UpdateRoutes(routes)
	}

	// use the retry method here because there is a chance the portal server isnt responding yet
	if err := util.Retry(updateRoute, 2, time.Second); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateRoutes(%+v): %s", routes, err.Error())
		return util.ErrorAppend(err, "failed to send routing updates to the router")
	}

	// update services
	services := generator.BuildServices(appModel)
	if err := client.UpdateServices(services); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateServices(%+v): %s", services, err.Error())
		return util.ErrorAppend(err, "failed to update port forwarding")
	}

	return nil
}

//
func portalClient(appModel *models.App) portal.PortalClient {
	return portal.New(appModel.LocalIPs["env"]+":8443", "123")
}
