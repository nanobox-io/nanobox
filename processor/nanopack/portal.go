package nanopack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nanobox-io/golang-portal-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type updatePortal struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("update_portal", updatePortalFunc)
}

func updatePortalFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return updatePortal{config}, nil
}

func (self updatePortal) Results() processor.ProcessConfig {
	return self.config
}

func (self updatePortal) Process() error {
	port := models.Service{}
	err := data.Get(util.AppName(), "portal", &port)
	if err != nil {
		return err
	}

	pClient := portal.New(port.ExternalIP+":8443", "123")

	// TODO: update portal
	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))

	services := []portal.Service{}
	routes := []portal.Route{}
	for _, node := range boxfile.Nodes("code") {
		service := models.Service{}
		err := data.Get(util.AppName(), node, &service)
		if err != nil {
			// unable to get the service
			continue
		}
		for _, service := range self.buildService(boxfile.Node(node), service) {
			if duplciateService(services, service) {
				if service.Port != 80 && service.Port != 443 {
					fmt.Println("duplicate port:", service.Port)
				}
			} else {
				services = append(services, service)
			}
		}
		for _, route := range self.buildRoutes(boxfile.Node(node), service) {
			if duplciateRoute(routes, route) {
				fmt.Println("duplicate route:", route.SubDomain, route.Path)
			} else {
				routes = append(routes, route)
			}

		}

	}

	// send to pulse
	err = pClient.UpdateServices(services)
	if err != nil {
		fmt.Println("update services", err)
		fmt.Printf("%+v\n", services)
		return err
	}

	err = pClient.UpdateRoutes(routes)
	if err != nil {
		fmt.Println("update routes", err)
		fmt.Printf("%+v\n", routes)
		return err
	}
	return nil
}

func duplciateService(services []portal.Service, service portal.Service) bool {
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
func duplciateRoute(services []portal.Route, service portal.Route) bool {
	for _, existingRoute := range services {
		if existingRoute.SubDomain == service.SubDomain && existingRoute.Path == service.Path {
			return true
		}
	}
	return false
}

// Server struct {
// 	// todo: change "Id" to "name" (for clarity)
// 	Id             string `json:"id,omitempty"`
// 	Host           string `json:"host"`
// 	Port           int    `json:"port"`
// 	Forwarder      string `json:"forwarder"`
// 	Weight         int    `json:"weight"`
// 	UpperThreshold int    `json:"upper_threshold"`
// 	LowerThreshold int    `json:"lower_threshold"`
// }
// Service struct {
// 	Id          string   `json:"id,omitempty"`
// 	Host        string   `json:"host"`
// 	Interface   string   `json:"interface,omitempty"`
// 	Port        int      `json:"port"`
// 	Type        string   `json:"type"`
// 	Scheduler   string   `json:"scheduler"`
// 	Persistence int      `json:"persistence"`
// 	Netmask     string   `json:"netmask"`
// 	Servers     []Server `json:"servers,omitempty"`
// }
func (self updatePortal) buildService(boxfile boxfile.Boxfile, service models.Service) []portal.Service {
	portServices := []portal.Service{}
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
			if portService.Type == "http" || portService.Type == "https" {
				portService.Servers[0].Host = "127.0.0.1"
				if portService.Type == "http" {
					portService.Servers[0].Port = 80
				} else {
					portService.Servers[0].Port = 443
				}
				portService.Type = "tcp"
			}
			portServices = append(portServices, portService)
		}
	}

	return portServices
}

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
func (self updatePortal) buildRoutes(boxfile boxfile.Boxfile, service models.Service) []portal.Route {
	portalRoutes := []portal.Route{}
	boxRoutes, ok := boxfile.Value("routes").([]string)
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
	for _, route := range boxRoutes {
		subdomain, path := parseRoute(route)
		portalRoute := portal.Route{
			SubDomain: subdomain,
			Path:      path,
		}

		for _, to := range ports(boxfile)["http"] {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, to))
		}
		for _, to := range ports(boxfile)["https"] {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, to))
		}
		if len(portalRoute.Targets) == 0 {
			portalRoute.Targets = append(portalRoute.Targets, fmt.Sprintf("http://%s:%s", service.InternalIP, "80"))
		}
		portalRoutes = append(portalRoutes, portalRoute)
	}

	return portalRoutes
}

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
	for _, port := range ports {
		p, ok := port.(string)
		if ok {
			portParts := strings.Split(p, ":")
			switch len(portParts) {
			case 1:
				rtn["http"][portParts[0]] = portParts[0]
			case 2:
				rtn["http"][portParts[0]] = portParts[1]
			case 3:
				switch portParts[0] {
				case "http", "https", "udp":
					rtn[portParts[0]][portParts[1]] = portParts[2]
				default:
					rtn["tcp"][portParts[1]] = portParts[2]
				}

			}
		}
		portInt, ok := port.(int)
		if ok {
			rtn["tcp"][strconv.Itoa(portInt)] = strconv.Itoa(portInt)
		}

	}
	return rtn
}
