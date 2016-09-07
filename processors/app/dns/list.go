package dns

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
)

// List lists all dns entries for an app
func List(a *models.App) error {

	// print the header
	fmt.Printf("\nDNS Aliases\n")
	
	// iterate
	for _, domain := range dns.List(a.ID) {
		fmt.Printf("  %s\n", domain.Domain)
	}
	
	fmt.Println()

	return nil
}
