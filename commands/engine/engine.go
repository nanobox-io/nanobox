// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package engine

import (
	"github.com/nanobox-io/nanobox-cli/auth"
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/s3"
	"github.com/spf13/cobra"
)

//
var (

	//
	EngineCmd = &cobra.Command{
		Use:   "engine",
		Short: "Subcommands to aid in developing a custom engine",
		Long:  ``,
	}

	//
	fFile  string
	Config = config.Default
	Auth   = auth.Default
	S3     = s3.Default
)

//
func init() {
	EngineCmd.AddCommand(fetchCmd)
	EngineCmd.AddCommand(newCmd)
	EngineCmd.AddCommand(publishCmd)
}
