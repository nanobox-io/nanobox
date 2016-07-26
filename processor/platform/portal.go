package platform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-portal-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processUpdatePortal ...
type processUpdatePortal struct {
	control processor.ProcessControl
	boxfile boxfile.Boxfile
	portal models.Service
}

//
func init() {
	processor.Register("update_portal", updatePortalFn)
}

//
func updatePortalFn(control processor.ProcessControl) (processor.Processor, error) {
	updatePortal := &processUpdatePortal{control: control}
	return updatePortal, updatePortal.ValidateMeta()
}

func (updatePortal *processUpdatePortal) ValidateMeta() error {
	if updatePortal.control.Meta["app_name"] == "" {
		updatePortal.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppID(), updatePortal.control.Env)
	}
	return nil
}

//
func (updatePortal processUpdatePortal) Results() processor.ProcessControl {
	return updatePortal.control
}

//
func (updatePortal *processUpdatePortal) Process() error {

	// load portal
	if err := updatePortal.loadPortal(); err != nil {
		return err
	}

	// load boxfile
	if err := updatePortal.loadBoxfile(); err != nil {
		return err
	}

	// update routes
	if err := updatePortal.updateRoutes(); err != nil {
		return err
	}

	// update ports
	return updatePortal.updatePorts()
}

// 
func (updatePortal *processUpdatePortal) loadPortal() error {
	return data.Get(updatePortal.control.Meta["app_name"], "portal", &updatePortal.portal)
}

// 
func (updatePortal *processUpdatePortal) loadBoxfile() error {
	updatePortal.boxfile = boxfile.New([]byte(updatePortal.control.Meta["boxfile"]))
	if !updatePortal.boxfile.Valid {
		return fmt.Errorf("invalid boxfile")
	}
	return nil
}

// update all the web routes that protal knows about
// updating the routes assumes the web servers are listening on
// 80 and 443 and in the container we assume the clients web server
// is listening on 8080
func (updatePortal *processUpdatePortal) updateRoutes() error {
	routes := []portal.Route{}

	// build the routes for all web containers
	for _, node := range updatePortal.boxfile.Nodes("web") {
		service := models.Service{}
		if err := data.Get(updatePortal.control.Meta["app_name"], node, &service); err != nil {
			continue // unable to get the service
		}

		for _, route := range updatePortal.buildRoutes(updatePortal.boxfile.Node(node), service) {
			lumber.Trace("route: %+v", route)
			if duplciateRoute(routes, route) {
				fmt.Println("duplicate route:", route.SubDomain, route.Path)
			}

			// append the new route to the routes we will register with portal
			routes = append(routes, route)
		}
	}

	// if i have a web and no routes i need to add a default one
	if len(updatePortal.boxfile.Nodes("web")) != 0 && len(routes) == 0 {
		webNode := updatePortal.boxfile.Nodes("web")[0]
		service := models.Service{}
		data.Get(updatePortal.control.Meta["app_name"], webNode, &service)
		routes = append(routes, portal.Route{
			Path:    "/",
			Targets: []string{fmt.Sprintf("http://%s:%s", service.InternalIP, "8080")},
		})
	}

	// send to portal
	lumber.Trace("new routes: %+v", routes)
	portalClient := portal.New(updatePortal.portal.ExternalIP+":8443", "123")
	return portalClient.UpdateRoutes(routes)
}

// Update the ports that portal knows about.
func (updatePortal *processUpdatePortal) updatePorts() error {
	services := []portal.Service{}

	//
	for _, node := range updatePortal.boxfile.Nodes("code") {
		service := models.Service{}
		if err := data.Get(updatePortal.control.Meta["app_name"], node, &service); err != nil {
			continue // unable to get the service
		}

		//
		for _, service := range updatePortal.buildService(updatePortal.boxfile.Node(node), service) {
			lumber.Trace("service: %+v", service)

			if duplicateService(services, service) {
				// if there is a duplicate port we will just contine and log
				fmt.Println("duplicate port: %+v", service.Port)
				continue
			}

			// add the new service to the list of services
			services = append(services, service)
		}
	}


	// send to portal
	lumber.Trace("new services: %+v", services)
	portalClient := portal.New(updatePortal.portal.ExternalIP+":8443", "123")
	return portalClient.UpdateServices(services)
}

// buildService builds all the tcp and udp port forwarding services
// it does not take into account any routing or information
func (updatePortal processUpdatePortal) buildService(boxfile boxfile.Boxfile, service models.Service) []portal.Service {

	portServices := []portal.Service{}

	//
	for protocol, protocolMap := range ports(boxfile) {
		for from, to := range protocolMap {
			fromInt, _ := strconv.Atoi(from)
			toInt, _ := strconv.Atoi(to)
			portService := portal.Service{
				Interface: "eth0",
				Port:      fromInt,
				Type:      protocol,
				Scheduler: "rr",
				Servers: []portal.Server{
					portal.Server{
						Host:      service.InternalIP,
						Port:      toInt,
						Forwarder: "m",
						Weight:    1,
					},
				},
			}

			portServices = append(portServices, portService)
		}
	}

	return portServices
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
func (updatePortal processUpdatePortal) buildRoutes(boxfile boxfile.Boxfile, service models.Service) []portal.Route {
	portalRoutes := []portal.Route{}
	boxRoutes, ok := boxfile.Value("routes").([]string)

	//
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

		portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, "8080"))
		portalRoutes = append(portalRoutes, portalRoute)
	}

	return portalRoutes
}

// duplicateService ...
func duplicateService(services []portal.Service, service portal.Service) bool {
	for _, existingService := range services {
		if existingService.Port == service.Port {
			return true
		}
	}
	return false
}

// duplicateRoute ...
func duplciateRoute(services []portal.Route, service portal.Route) bool {
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

// ports ...
func ports(box boxfile.Boxfile) map[string]map[string]string {
	// we allow tcp and udp ports
	rtn := map[string]map[string]string{
		"tcp":   map[string]string{},
		"udp":   map[string]string{},
	}

	// get the boxfiles ports section
	ports, ok := box.Value("ports").([]interface{})
	if !ok {
		return rtn
	}

	// loop through the given ports and create hash data
	// for each one.
	for _, port := range ports {
		p, ok := port.(string)
		if ok {
			portParts := strings.Split(p, ":")
			switch len(portParts) {
			case 1:
				rtn[TCP][portParts[0]] = portParts[0]
			case 2:
				rtn[TCP][portParts[0]] = portParts[1]
			case 3:
				// the first part needs to be tcp or udp
				// if it is neither we just assume tcp
				switch portParts[0] {
				case UDP:
					rtn[portParts[0]][portParts[1]] = portParts[2]
				default:
					rtn[TCP][portParts[1]] = portParts[2]
				}

			}
		}
		// if only a number is provided we assume tcp:num:num
		portInt, ok := port.(int)
		if ok {
			rtn[TCP][strconv.Itoa(portInt)] = strconv.Itoa(portInt)
		}

	}

	return rtn
}
