//
package server

import ()

// Update issues an update to nanobox server
func Update(params string) error {

	if _, err := Post("/image-update?"+params, "text/plain", nil); err != nil {
		return err
	}

	return nil
}
