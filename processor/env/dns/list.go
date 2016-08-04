package dns

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dns"
)

// processEnvDNSList ...
type processEnvDNSList struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("env_dns_list", envDNSListFn)
}

// envDNSListFn creates a processEnvDNSList and validates the meta data in the control
func envDNSListFn(control processor.ProcessControl) (processor.Processor, error) {
	envDNSList := &processEnvDNSList{control: control}
	return envDNSList, envDNSList.validateMeta()
}

func (envDNSList processEnvDNSList) Results() processor.ProcessControl {
	return envDNSList.control
}

//
func (envDNSList processEnvDNSList) Process() error {

	appID := fmt.Sprintf("%s_%s", config.AppID(), envDNSList.control.Env)

	fmt.Printf("dns entries for %s(%s):\n", config.AppName(), envDNSList.control.Env)
	entries := dns.List(appID)
	for _, entry := range entries {
		fmt.Printf("  %s\n", entry.Domain)
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (envDNSList *processEnvDNSList) validateMeta() error {
	// currently does nothing but may in the future
	return nil
}
