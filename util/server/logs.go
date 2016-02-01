//
package server

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
)

// Logs diplayes historical logs from the server
func Logs(params string) error {

	logs := []mistutil.Log{}

	//
	if _, err := Get("/logs?"+params, &logs); err != nil {
		return err
	}

	//
	fmt.Printf(stylish.Bullet("Showing last %v entries:", len(logs)))
	for _, log := range logs {
		mistutil.ProcessLog(log)
	}

	return nil
}
