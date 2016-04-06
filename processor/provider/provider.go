package provider

import (
	"github.com/nanobox-io/nanobox/processor"
)

func init() {
	processor.Register("provider_setup", providerSetupFunc)
	processor.Register("provider_destroy", providerDestroyFunc)
}
