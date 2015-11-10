//
package server

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/util/server/mist"
)

// Logs diplayes historical logs from the server
func Logs(params string) error {

	logs := []mist.Log{}

	//
	res, err := Get("/logs?"+params, &logs)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//
	fmt.Printf(stylish.Bullet("Showing last %v entries:", len(logs)))
	for _, log := range logs {
		mist.ProcessLog(log)
	}

	return nil
}
