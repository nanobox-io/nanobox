package router

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/golang-portal-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
)

func BuildRoutes(appModel *models.App) []portal.Route {
	boxfile := loadBoxfile(appModel)
	routes := []portal.Route{}

	// build the routes for all web containers
	for _, node := range boxfile.Nodes("web") {
		component, err := models.FindComponentBySlug(appModel.ID, node)
		if err != nil {
			continue // unable to get the component
		}

		for _, route := range buildComponentRoutes(boxfile.Node(node), component) {
			if duplicateRoute(routes, route) {
				continue // this route exits already so we wont replace it
			}

			// append the new route to the routes we will register with portal
			routes = append(routes, route)
		}
	}

	// if i have a web and no routes i need to add a default one
	if len(boxfile.Nodes("web")) != 0 && len(routes) == 0 {
		webNode := boxfile.Nodes("web")[0]
		component, _ := models.FindComponentBySlug(appModel.ID, webNode)

		routes = append(routes, portal.Route{
			Path:    "/",
			Targets: []string{fmt.Sprintf("http://%s:%s", component.IPAddr(), "8080")},
		})
	}

	// send to portal
	return routes
}

// buildRoutes ...
//
// Route struct {
// 	// defines match characteristics
// 	SubDomain string `json:"subdomain"` // subdomain to match on - "admin"
// 	Domain    string `json:"domain"`    // domain to match on - "myapp.com"
// 	Path      string `json:"path"`      // route to match on - "/admin"
// 	// defines actions
// 	Targets []string `json:"targets"` // ips of servers - ["http://127.0.0.1:8080/app1","http://127.0.0.2"] (optional)
// 	FwdPath string   `json:"fwdpath"` // path to forward to targets - "/goadmin" incoming req: test.com/admin -> 127.0.0.1/goadmin (optional)
// 	Page    string   `json:"page"`    // page to serve instead of routing to targets - "<HTML>We are fixing it</HTML>" (optional)
// }
func buildComponentRoutes(boxfile boxfile.Boxfile, component *models.Component) []portal.Route {
	portalRoutes := []portal.Route{}
	boxRoutes, ok := boxfile.Value("routes").([]string)

	// if the routes are not a []strings try converting
	// an []interface to []string
	if !ok {
		tmps, ok := boxfile.Value("routes").([]interface{})
		if !ok {
			// no routes apparently
			return portalRoutes
		}
		for _, tmp := range tmps {
			if str, ok := tmp.(string); ok {
				boxRoutes = append(boxRoutes, str)
			}
		}
	}

	//
	for _, route := range boxRoutes {
		subdomain, path := parseRoute(route)
		portalRoute := portal.Route{
			SubDomain: subdomain,
			Path:      path,
		}

		portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", component.IPAddr(), "8080"))
		portalRoutes = append(portalRoutes, portalRoute)
	}

	return portalRoutes
}

// duplicateRoute ...
func duplicateRoute(services []portal.Route, service portal.Route) bool {
	for _, existingRoute := range services {
		if existingRoute.SubDomain == service.SubDomain && existingRoute.Path == service.Path {
			return true
		}
	}
	return false
}

// parseRoute ...
func parseRoute(route string) (subdomain, path string) {
	routeParts := strings.Split(route, ":")
	switch len(routeParts) {
	case 1:
		path = routeParts[0]
	case 2:
		subdomain = routeParts[0]
		path = routeParts[1]
	}
	return
}
