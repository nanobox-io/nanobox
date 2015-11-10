//
package vagrant

// Install downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Install() error {

	// download nanobox image
	if err := download(); err != nil {
		return err
	}

	// add nanobox image
	return add()
}
