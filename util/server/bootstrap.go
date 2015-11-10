//
package server

import ()

// Bootstrap issues a bootstrap to nanobox server
func Bootstrap(params string) error {

	res, err := Post("/bootstrap?"+params, "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
