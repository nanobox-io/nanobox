// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/file"
	"github.com/nanobox-io/nanobox/util/file/hosts"
	"github.com/nanobox-io/nanobox/util/notify"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/server"
	"github.com/nanobox-io/nanobox/util/server/mist"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

// These public variables are for replacement testing
var (
	Config  = config.Default
	Util    = util.Default
	Server  = server.Default
	Mist    = mist.Default
	Vagrant = vagrant.Default
	Hosts   = hosts.Default
	File    = file.Default
	Print   = print.Default
	Notify  = notify.Default
)
