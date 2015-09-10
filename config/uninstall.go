// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

// import (
// 	"fmt"
// 	"os"
//
// 	"github.com/pagodabox/nanobox-cli/util"
// )
//
// //
// func Uninstall(force bool) {
//
// 	//
// 	if force != true {
//
// 		response := util.Prompt("Are you sure you want to uninstall the Pagoda Box CLI (y/N)? ")
//
// 		if response != "y" {
// 			fmt.Printf("'%v' - Pagoda Box CLI will not be uninstalled. Exiting...\n", response)
// 			os.Exit(0)
// 		}
// 	}
//
// 	fmt.Print("Uninstalling... ")
//
// 	//
// 	if err := os.RemoveAll(nanoDir); err != nil {
// 		util.LogFatal("[install] os.Remove() failed", err)
// 	}
//
// 	util.CPrint("[green]success[reset]")
// }
