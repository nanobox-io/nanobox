package platform

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-portal-client"
	generator "github.com/nanobox-io/nanobox/generators/router"
	"github.com/nanobox-io/nanobox/models"
)

// UpdatePortal ...
func UpdatePortal(appModel *models.App) error {
	client := portalClient(appModel)

	// update routes
	routes := generator.BuildRoutes(appModel)
	if err := client.UpdateRoutes(routes); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateRoutes(%+v): %s", routes, err.Error())
		return fmt.Errorf("failed to sending routing updates to the router: %s", err.Error())
	}

	// update services
	services := generator.BuildServices(appModel)
	if err := client.UpdateServices(services); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateServices(%+v): %s", services, err.Error())
		return fmt.Errorf("failed to update port forwarding: %s", err.Error())
	}

	return nil
}

//
func portalClient(appModel *models.App) portal.PortalClient {
	return portal.New(appModel.GlobalIPs["env"]+":8443", "123")
}
