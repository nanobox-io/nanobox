//
package server

import (
	"net/http"
	"time"

	"github.com/nanobox-io/nanobox/config"
)

// Ping issues a ping to nanobox server
func Ping() (bool, error) {

	// a new client is used to allow for shortening the request timeout
	client := http.Client{Timeout: time.Duration(2 * time.Second)}

	//
	res, err := client.Get(config.ServerURL + "/ping")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	//
	return res.StatusCode/100 == 2, nil
}
