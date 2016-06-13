package platform

import (
	"fmt"
	"strconv"
	"strings"

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
}

//
func init() {
	processor.Register("update_portal", updatePortalFunc)
}

//
func updatePortalFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return processUpdatePortal{control}, nil
}

//
func (updatePortal processUpdatePortal) Results() processor.ProcessControl {
	return updatePortal.control
}

//
func (updatePortal processUpdatePortal) Process() error {
	port := models.Service{}

	//
	if err := data.Get(config.AppName(), "portal", &port); err != nil {
		return err
	}

	pClient := portal.New(port.ExternalIP+":8443", "123")
	boxfile := boxfile.New([]byte(updatePortal.control.Meta["boxfile"]))
	services := []portal.Service{}
	routes := []portal.Route{}

	//
	for _, node := range boxfile.Nodes("code") {
		service := models.Service{}
		if err := data.Get(config.AppName(), node, &service); err != nil {
			continue // unable to get the service
		}

		//
		for _, service := range updatePortal.buildService(boxfile.Node(node), service) {
			if duplicateService(services, service) {
				if service.Port != 80 && service.Port != 443 {
					fmt.Println("duplicate port:", service.Port)
				}
			} else {
				services = append(services, service)
			}
		}

		//
		for _, route := range updatePortal.buildRoutes(boxfile.Node(node), service) {
			if duplciateRoute(routes, route) {
				fmt.Println("duplicate route:", route.SubDomain, route.Path)
			} else {
				routes = append(routes, route)
			}
		}
	}

	// if i have a web and no services i need to add a default one
	if len(boxfile.Nodes("web")) != 0 && len(services) == 0 {

		//
		services = append(services, portal.Service{
			Interface: "eth0",
			Port:      80,
			Type:      TCP,
			Scheduler: "rr",
			Servers: []portal.Server{
				portal.Server{
					Host:      "127.0.0.1",
					Port:      80,
					Forwarder: "m",
					Weight:    1,
				},
			},
		})

		//
		services = append(services, portal.Service{
			Interface: "eth0",
			Port:      443,
			Type:      TCP,
			Scheduler: "rr",
			Servers: []portal.Server{
				portal.Server{
					Host:      "127.0.0.1",
					Port:      443,
					Forwarder: "m",
					Weight:    1,
				},
			},
		})
	}

	// if i have a web and no routes i need to add a default one
	if len(boxfile.Nodes("web")) != 0 && len(routes) == 0 {
		webNode := boxfile.Nodes("web")[0]
		service := models.Service{}
		data.Get(config.AppName(), webNode, &service)
		routes = append(routes, portal.Route{
			Path:    "/",
			Targets: []string{fmt.Sprintf("http://%s:%s", service.InternalIP, "80")},
		})
	}

	// send to pulse
	if err := pClient.UpdateServices(services); err != nil {
		return err
	}

	//
	if err := pClient.UpdateRoutes(routes); err != nil {
		return err
	}

	return nil
}

// duplicateService ...
func duplicateService(services []portal.Service, service portal.Service) bool {
	if service.Port == 80 || service.Port == 443 {
		return false
	}
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

// buildService ...
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

			//
			if portService.Type == HTTP || portService.Type == HTTPS {
				portService.Servers[0].Host = "127.0.0.1"
				if portService.Type == HTTP {
					portService.Servers[0].Port = 80
				} else {
					portService.Servers[0].Port = 443
				}
				portService.Type = TCP
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

		for _, to := range ports(boxfile)[HTTP] {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, to))
		}
		for _, to := range ports(boxfile)[HTTPS] {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, to))
		}
		if len(portalRoute.Targets) == 0 {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, "80"))
		}
		portalRoutes = append(portalRoutes, portalRoute)
	}

	return portalRoutes
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
	rtn := map[string]map[string]string{
		"http":  map[string]string{},
		"https": map[string]string{},
		"tcp":   map[string]string{},
		"udp":   map[string]string{},
	}

	ports, ok := box.Value("ports").([]interface{})
	if !ok {
		return rtn
	}

	//
	for _, port := range ports {
		p, ok := port.(string)
		if ok {
			portParts := strings.Split(p, ":")
			switch len(portParts) {
			case 1:
				rtn[HTTP][portParts[0]] = portParts[0]
			case 2:
				rtn[HTTP][portParts[0]] = portParts[1]
			case 3:
				switch portParts[0] {
				case HTTP, HTTPS, UDP:
					rtn[portParts[0]][portParts[1]] = portParts[2]
				default:
					rtn[TCP][portParts[1]] = portParts[2]
				}

			}
		}
		portInt, ok := port.(int)
		if ok {
			rtn[TCP][strconv.Itoa(portInt)] = strconv.Itoa(portInt)
		}

	}

	return rtn
}
