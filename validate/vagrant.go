package validate

import (
	"fmt"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

func init() {
	Register("vagrant", vagrantCheck)
}

func vagrantCheck() error {
	if exists := vagrant.Exists(); !exists {
		return fmt.Errorf("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
	}	
	return nil
}
