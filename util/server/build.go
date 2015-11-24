//
package server

import ()

// Build issues a build to nanobox server
func Build(params string) error {

	if _, err := Post("/builds?"+params, "text/plain", nil); err != nil {
		return err
	}

	return nil
}
