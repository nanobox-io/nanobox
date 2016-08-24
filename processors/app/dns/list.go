package dns

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
)

// List lists all dns entries for an app
func List(a *models.App) error {

	fmt.Printf("dns entries for %s(%s):\n", a.EnvID, a.Name)
	entries := dns.List(a.ID)
	for _, entry := range entries {
		fmt.Printf("  %s\n", entry.Domain)
	}

	return nil
}
