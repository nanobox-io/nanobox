// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"

	"github.com/spf13/cobra"
)

//
var nanoboxPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "(coming soon)",
	Long: `
Description:
  Pushed app to production`,

	Run: nanoboxPublish,
}

// nanoboxPublish
func nanoboxPublish(ccmd *cobra.Command, args []string) {
	fmt.Println("coming soon (http://production.nanobox.io)")
}
