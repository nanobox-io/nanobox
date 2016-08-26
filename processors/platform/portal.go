package platform

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-portal-client"

	route_generator "github.com/nanobox-io/nanobox/generators/router"
	"github.com/nanobox-io/nanobox/models"
)

//
func UpdatePortal(appModel *models.App) error {

	// update routes
	routes := route_generator.BuildRoutes(appModel)
	if err := portalClient(appModel).UpdateRoutes(routes); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateRoutes(%+v): %s", routes, err.Error)
		return fmt.Errorf("failed to sending routing updates to the router: %s",err.Error())
	}

	// update services
	services := route_generator.BuildServices(appModel)
	if err := portalClient(appModel).UpdateServices(services); err != nil {
		lumber.Error("platform:UpdatePortal:UpdateServices(%+v): %s", services, err.Error)
		return fmt.Errorf("failed to update port forwarding: %s", err.Error())
	}

	return nil
}


//
func portalClient(appModel *models.App) portal.PortalClient {
	return portal.New(appModel.GlobalIPs["env"]+":8443", "123")
}
