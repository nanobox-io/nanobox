//
package server

import ()

// Build issues a build to nanobox server
func Build(params string) error {

	res, err := Post("/builds?"+params, "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
