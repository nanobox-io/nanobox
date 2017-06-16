package router

import (
	"strconv"
	"strings"

	"github.com/nanobox-io/golang-portal-client"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
)

// Update the ports that portal knows about.
func BuildServices(appModel *models.App) []portal.Service {
	services := []portal.Service{}
	boxfile := loadBoxfile(appModel)
	//
	for _, node := range boxfile.Nodes("code") {
		component, err := models.FindComponentBySlug(appModel.ID, node)
		if err != nil {
			continue // unable to get the component
		}

		//
		for _, service := range buildComponentServices(boxfile.Node(node), component) {

			if duplicateService(services, service) {
				continue // if there is a duplicate port we will just contine
			}

			// add the new service to the list of services
			services = append(services, service)
		}
	}

	// send to portal
	return services
}

// buildService builds all the tcp and udp port forwarding services
// it does not take into account any routing or information
func buildComponentServices(boxfile boxfile.Boxfile, component *models.Component) []portal.Service {

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
					{
						Host:      component.IPAddr(),
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

// duplicateService ...
func duplicateService(services []portal.Service, service portal.Service) bool {
	for _, existingService := range services {
		if existingService.Port == service.Port {
			return true
		}
	}
	return false
}

// ports ...
func ports(box boxfile.Boxfile) map[string]map[string]string {
	// we allow tcp and udp ports
	rtn := map[string]map[string]string{
		"tcp": {},
		"udp": {},
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
				rtn["tcp"][portParts[0]] = portParts[0]
			case 2:
				rtn["tcp"][portParts[0]] = portParts[1]
			case 3:
				// the first part needs to be tcp or udp
				// if it is neither we just assume tcp
				switch portParts[0] {
				case "udp":
					rtn[portParts[0]][portParts[1]] = portParts[2]
				case "tcp":
					rtn["tcp"][portParts[1]] = portParts[2]
				default:
					display.BadPortType(portParts[0])
					rtn["tcp"][portParts[1]] = portParts[2]
				}

			}
		}
		// if only a number is provided we assume tcp:num:num
		portInt, ok := port.(int)
		if ok {
			rtn["tcp"][strconv.Itoa(portInt)] = strconv.Itoa(portInt)
		}

	}

	return rtn
}
