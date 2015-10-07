// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package box

import "github.com/spf13/cobra"

//
var BoxCmd = &cobra.Command{
	Use:   "box",
	Short: "",
	Long:  ``,
}

//
func init() {
	BoxCmd.AddCommand(installCmd)
	BoxCmd.AddCommand(updateCmd)
}
