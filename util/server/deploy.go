//
package server

import ()

// Deploy issues a deploy to nanobox server
func Deploy(params string) error {

	res, err := Post("/deploys?"+params, "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
