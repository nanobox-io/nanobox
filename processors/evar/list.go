package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

func List(appModel *models.App) error {

	// print the header
	fmt.Printf("\nEnvironment Variables\n")

	// iterate
	for key, val := range appModel.Evars {
		fmt.Printf("  %s = %s\n", key, val)
	}

	fmt.Println()

	return nil
}
