// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package engine

import "github.com/spf13/cobra"

//
var (

	//
	EngineCmd = &cobra.Command{
		Use:   "engine",
		Short: "",
		Long:  ``,
	}

	//
	fFile string
)

//
func init() {
	EngineCmd.AddCommand(fetchCmd)
	EngineCmd.AddCommand(newCmd)
	EngineCmd.AddCommand(publishCmd)
}
