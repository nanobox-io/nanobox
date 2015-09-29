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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var boxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: boxUpdate,
}

//
func getMD5() string {
	md5, err := http.Get("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5")
	if err != nil {
		config.Fatal("[commands/boxInstall] http.Get() failed", err.Error())
	}
	defer md5.Body.Close()

	b, err := ioutil.ReadAll(md5.Body)
	if err != nil {
		config.Fatal("[commands/boxInstall] ioutil.ReadAll() failed", err.Error())
	}

	return string(b)
}

//
func compairMD5() string {
	f, err := os.Open(config.Root + "/nanobox-boot2docker.md5")
	if err != nil {
		return ""
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		config.Fatal("[commands/boxInstall] ioutil.ReadAll() failed", err.Error())
	}

	return string(b)
}

// boxUpdate
func boxUpdate(ccmd *cobra.Command, args []string) {

	// check to see if an update is even needed
	if getMD5() == compairMD5() {
		return
	}

	//
	// if !needBox() {
	// 	fmt.Printf(stylish.Bullet("Uninstalling previous virtual machine image..."))
	// 	if err := exec.Command("vagrant", "box", "remove", "--force", "nanobox/boot2docker").Run(); err != nil {
	// 		config.Fatal("[commands/boxUpdate] exec.Command() failed", err.Error())
	// 	}
	// }

	//
	boxInstall(nil, args)
}
