package code

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/print"
)

// processCodeSetup ...
type processCodeSetup struct {
	control processor.ProcessControl
	service models.Service
	fail    bool
}

//
func init() {
	processor.Register("code_setup", codeSetupFn)
}

//
func codeSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return &processCodeSetup{control: control}, nil
}

//
func (codeSetup *processCodeSetup) clean(fn func()) func() {
	return func() {
		if codeSetup.fail {
			fn()
		}
	}
}

//
func (codeSetup processCodeSetup) Results() processor.ProcessControl {
	return codeSetup.control
}

//
func (codeSetup *processCodeSetup) Process() error {

	// make sure i was given a name and image
	if codeSetup.control.Meta["name"] == "" || codeSetup.control.Meta["image"] == "" {
		return errMissingImageOrName
	}

	// quit early if the service was found in the database as well as docker
	if codeSetup.serviceExists() {
		return nil
	}

	//
	if err := codeSetup.getLocalIP(); err != nil {
		return err
	}
	defer codeSetup.clean(codeSetup.returnLocalIP)()

	//
	if err := codeSetup.getGlobalIP(); err != nil {
		return err
	}
	defer codeSetup.clean(codeSetup.returnGlobalIP)()

	// pull the docker image
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(codeSetup.control.DisplayLevel), codeSetup.control.Meta["image"])
	if _, err := docker.ImagePull(codeSetup.control.Meta["image"], &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}

	//
	if err := codeSetup.createContainer(); err != nil {
		return err
	}
	defer codeSetup.clean(codeSetup.removeContainer)()

	//
	if err := codeSetup.addIPToProvider(); err != nil {
		return err
	}
	defer codeSetup.clean(codeSetup.removeIPFromProvider)()

	// save the service
	if err := data.Put(config.AppName(), codeSetup.control.Meta["name"], codeSetup.service); err != nil {
		codeSetup.fail = true
		lumber.Error("insert data: ", err)
		return err
	}

	return nil
}

// serviceExists ...
func (codeSetup *processCodeSetup) serviceExists() bool {
	service := models.Service{}
	databaseErr := data.Get(config.AppName(), codeSetup.control.Meta["name"], &service)

	// set the service i found so i dont re allocate ips
	if databaseErr == nil {
		codeSetup.service = service
	}
	_, containerErr := docker.GetContainer(service.ID)

	return databaseErr == nil && containerErr == nil
}

// getLocalIP ...
func (codeSetup *processCodeSetup) getLocalIP() error {

	// if the service already has an ip
	if codeSetup.service.InternalIP != "" {
		return nil
	}
	ip, err := dhcp.ReserveLocal()
	if err != nil {
		codeSetup.fail = true
		return err
	}
	codeSetup.service.InternalIP = ip.String()

	return nil
}

// getGlobalIP ...
func (codeSetup *processCodeSetup) getGlobalIP() error {
	if codeSetup.service.ExternalIP != "" {
		// if the service already has an ip
		return nil
	}
	ip, err := dhcp.ReserveGlobal()
	if err != nil {
		codeSetup.fail = true
		return err
	}
	codeSetup.service.ExternalIP = ip.String()

	return nil
}

// returnLocalIP ...
func (codeSetup *processCodeSetup) returnLocalIP() {
	dhcp.ReturnIP(net.ParseIP(codeSetup.service.InternalIP))
}

// returnGlobalIP ...
func (codeSetup *processCodeSetup) returnGlobalIP() {
	dhcp.ReturnIP(net.ParseIP(codeSetup.service.ExternalIP))
}

// addIPToProvider ...
func (codeSetup *processCodeSetup) addIPToProvider() error {
	if codeSetup.service.InternalIP == "" || codeSetup.service.ExternalIP == "" {
		return fmt.Errorf("im missing an ip to bind to the provider")
	}

	if err := provider.AddIP(codeSetup.service.ExternalIP); err != nil {
		codeSetup.fail = true
		return err
	}

	if err := provider.AddNat(codeSetup.service.ExternalIP, codeSetup.service.InternalIP); err != nil {
		codeSetup.fail = true
		return err
	}
	return nil
}

// removeIPFromProvider ...
func (codeSetup *processCodeSetup) removeIPFromProvider() {
	provider.RemoveNat(codeSetup.service.ExternalIP, codeSetup.service.InternalIP)
	provider.RemoveIP(codeSetup.service.ExternalIP)
}

// createContainer ...
func (codeSetup *processCodeSetup) createContainer() error {
	// configure the container
	fmt.Println("-> building container", codeSetup.control.Meta["name"])
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-%s-%s", config.AppName(), codeSetup.control.Env, codeSetup.control.Meta["name"]),
		Image:   codeSetup.control.Meta["image"],
		Network: "virt",
		IP:      codeSetup.service.InternalIP,
	}

	// create docker container
	container, err := docker.CreateContainer(config)
	if err != nil {
		codeSetup.fail = true
		lumber.Error("container: ", err)
		return err
	}
	codeSetup.service.ID = container.ID
	codeSetup.service.Name = codeSetup.control.Meta["name"]
	codeSetup.service.Type = "code"
	return nil
}

// removeContainer ...
func (codeSetup *processCodeSetup) removeContainer() {
	// catch the error here and display it but dont error out
	docker.ContainerRemove(codeSetup.service.ID)
}
