//
package server

import ()

// Deploy issues a deploy to nanobox server
func Deploy(params string) error {

	if _, err := Post("/deploys?"+params, "text/plain", nil); err != nil {
		return err
	}

	return nil
}
