// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// UpgradeCommand satisfies the Command interface for obtaining user info
type UpgradeCommand struct{}

// Help
func (c *UpgradeCommand) Help() {
	ui.CPrint(`
Description:
  Updates the nanobox docker images

Usage:
  pagoda upgrade
  `)
}

// Run
func (c *UpgradeCommand) Run(opts []string) {

	//
	fmt.Printf(stylish.Bullet("Updating nanobox docker images..."))
	if err := api.DoRawRequest(nil, "POST", fmt.Sprintf("http://%v:1757/image-update", config.Nanofile.IP), nil, nil); err != nil {
		fmt.Printf(stylish.Error("failed to update nanobox docker images", fmt.Sprintf("nanobox was unable to updated its docker images, and failed with the following error: %v", err)))
	}
	fmt.Printf(stylish.Success())
}
