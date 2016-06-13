// Package provider ...
package provider

import "github.com/nanobox-io/nanobox/processor"

//
func init() {
	processor.Register("provider_setup", providerSetupFunc)
	processor.Register("provider_destroy", providerDestroyFunc)
	processor.Register("provider_stop", providerStopFunc)
}
