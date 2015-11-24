//
package server

import ()

// Bootstrap issues a bootstrap to nanobox server
func Bootstrap(params string) error {

	if _, err := Post("/bootstrap?"+params, "text/plain", nil); err != nil {
		return err
	}

	return nil
}
