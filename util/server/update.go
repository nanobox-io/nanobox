//
package server

import ()

// Update issues an update to nanobox server
func Update(params string) error {

	res, err := Post("/image-update?"+params, "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
