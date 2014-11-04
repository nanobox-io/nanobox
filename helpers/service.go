package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/ui"
)

// GetServiceBySlug attemtps to find an app service by name or UID. It takes an
// app name, fetches all its services, then iterates over the services to see if
// there is a name or UID that matches the provided slug.
func GetServiceBySlug(appName, serviceSlug string, api *pagodaAPI.Client) (*pagodaAPI.AppService, error) {

	services, err := api.GetAppServices(appName)
	if err != nil {
		return nil, err
	}

	//
	for _, s := range services {
		if s.Name == serviceSlug || s.UID == serviceSlug {
			service, err := api.GetAppService(appName, s.ID)
			if err != nil {
				return nil, err
			}

			return service, nil
		}
	}

	return nil, errors.New("We couldn't find a service matching the name '" + serviceSlug + "' for the app '" + appName + "'.")
}

// EnablePublicTunnel enables a service's public tunnel if an SSH command is run
// and the public tunnel for that service has not yet been enabled
func EnablePublicTunnel(service *pagodaAPI.AppService, api *pagodaAPI.Client, opts *SSHOptions) {

	appServiceUpdateOptions := &pagodaAPI.AppServiceUpdateOptions{PublicTunnel: true}

	service, err := api.UpdateAppService(service.AppID, service.ID, appServiceUpdateOptions)
	if err != nil {
		fmt.Println("There was a problem enabling '%s' SSH. See ~/.pagodabox/log.txt for details", service.UID)
		ui.Error("pagoda helpers.EnablePublicTunnel", err)
	}

	fmt.Print("Enabling tunnel..")

	for {
		fmt.Print(".")

		service, err := api.GetAppService(service.AppID, service.ID)
		if err != nil {
			fmt.Printf("Oops! We could not find a '%s'.\n", service.UID)
			os.Exit(1)
		}

		if service.TunnelIP != "" && service.TunnelPort != 0 {
			opts.RemoteIP = service.TunnelIP
			opts.RemotePort = service.TunnelPort

			fmt.Println("")

			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

// DetermineAppStatus taks the string 'status' of a services and returns a color
// representing that status
func DetermineServiceStatus(s string) string {
	switch s {

	//
	case "initialized":
		return "[blue]"

	//
	case "active":
		return "[green]"

	//
	case "inactive", "defunct":
		return "[red]"
	}

	return ""
}
