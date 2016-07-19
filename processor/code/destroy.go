package code

import (
	"fmt"
	"net"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// processCodeDestroy ...
type processCodeDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("code_destroy", codeDestroyFn)
}

//
func codeDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	if control.Meta["name"] == "" {
		return nil, errMissingImageOrName
	}
	return &processCodeDestroy{control: control}, nil
}

//
func (codeDestroy processCodeDestroy) Results() processor.ProcessControl {
	return codeDestroy.control
}

//
func (codeDestroy *processCodeDestroy) Process() error {

	// get the service from the database
	service := models.Service{}

	//
	bucket := fmt.Sprintf("%s_%s", config.AppName(), codeDestroy.control.Env)
	if err := data.Get(bucket, codeDestroy.control.Meta["name"], &service); err != nil {
		return err
	}

	//
	if err := docker.ContainerRemove(service.ID); err != nil {
		return err
	}

	//
	if err := provider.RemoveNat(service.ExternalIP, service.InternalIP); err != nil {
		return err
	}

	//
	if err := provider.RemoveIP(service.ExternalIP); err != nil {
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(service.ExternalIP)); err != nil {
		return err
	}

	//
	if err := dhcp.ReturnIP(net.ParseIP(service.InternalIP)); err != nil {
		return err
	}

	// remove the service from the database
	return data.Delete(bucket, codeDestroy.control.Meta["name"])
}
