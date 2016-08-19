package dns

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
)

// List ...
type List struct {
	App models.App
}

//
func (list List) Run() error {

	fmt.Printf("dns entries for %s(%s):\n", list.App.EnvID, list.App.Name)
	entries := dns.List(list.App.ID)
	for _, entry := range entries {
		fmt.Printf("  %s\n", entry.Domain)
	}

	return nil
}
