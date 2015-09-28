// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "strings"

	// semver "github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"

	// "github.com/pagodabox/nanobox-cli/config"
	// "github.com/pagodabox/nanobox-cli/util"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

type asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

type release struct {
	Assets  []asset
	Version string `json:"name"`
}

//
var boxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: boxInstall,
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {
	releases := []release{}

	//
	res, err := http.Get("https://api.github.com/repos/pagodabox/nanobox-boot2docker/releases")
	if err != nil {
		fmt.Println("BRANK!", err)
	}
	defer res.Body.Close()

	//
	body, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &releases); err != nil {
		fmt.Println("BRUNK!", err)
	}

	fmt.Println("BODY! %#v\n", string(body))

	// latestRelease := release{}
	// latest, err := semver.NewVersion("0.0.0")
	// if err != nil {
	// 	fmt.Println("BUZUNK!", err)
	// }

	// //
	// for _, release := range releases {
	// 	if release.Name != "" {
	//
	// 		ver, err := semver.NewVersion(strings.Replace(release.Name, "v", "", -1))
	// 		if err != nil {
	// 			fmt.Println("BRUOONK!", err)
	// 		}
	//
	// 		if latest.LessThan(*ver) {
	// 			latestRelease = release
	// 			latest = ver
	// 		}
	// 	}
	// }

	// fmt.Printf("LATEST RELEASE! %#v\n", latestRelease)

	// if currentVersion() != 0 {
	// 	fmt.Println("I already have a version")
	// 	return
	// }
	// release := latestVersion()
	// asset := release.Assets[0]
	// put file downloader here downloading from asset.DownloadURL
	// setVersion(release.version())
	// vagrant box add ~/.nanobox/boot2docker.box
}

// currentVersion :=
//
// func currentVersion() int {
// 	ver, err := ioutil.ReadFile("/Users/sdomino/.nanobox/.box")
// 	if err != nil {
// 		return 0
// 	}
// 	verInt, _ := strconv.Atoi(string(ver))
// 	return verInt
// }

// func releases() []release {
// 	releases := []release{}
// 	resp, err := http.Get("https://api.github.com/repos/pagodabox/nanobox-boot2docker/releases")
// 	if err != nil {
// 		return releases
// 	}
//
// 	body, err := ioutil.ReadAll(resp.Body)
// 	err = json.Unmarshal(body, &releases)
// 	if err != nil {
// 		return releases
// 	}
// 	return releases
// }
