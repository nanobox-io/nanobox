// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	// "bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "os/exec"
	// "strings"

	"github.com/nanobox-io/nanobox-cli/config"
)

//
// func NeedBox() bool {
//
// 	need := true
//
// 	cmd := exec.Command("vagrant", "box", "list")
//
// 	//
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		config.Fatal("[commands/boxInstall] cmd.StdoutPipe() failed", err.Error())
// 	}
//
// 	//
// 	scanner := bufio.NewScanner(stdout)
// 	go func() {
// 		for scanner.Scan() {
// 			if strings.HasPrefix(scanner.Text(), "nanobox/boot2docker") {
// 				need = false
// 			}
// 		}
// 	}()
//
// 	//
// 	if err := cmd.Start(); err != nil {
// 		config.Fatal("[commands/boxInstall] cmd.Start() failed", err.Error())
// 	}
//
// 	//
// 	if err := cmd.Wait(); err != nil {
// 		config.Fatal("[commands/boxInstall] cmd.Wait() failed", err.Error())
// 	}
//
// 	return need
// }

// VMRemoteMD5
func VMRemoteMD5() string {
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

// VMLocalMD5
func VMLocalMD5() string {
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

// VMDownload
func VMDownload() {

	//
	box, err := os.Create(config.Root + "/nanobox-boot2docker.box")
	if err != nil {
		config.Fatal("[commands/update] os.Create() failed", err.Error())
	}
	defer box.Close()

	// download the box with a progres bar
	Progress(fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.box"), box)
}
