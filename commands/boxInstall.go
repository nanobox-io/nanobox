// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	semver "github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
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

func hasBox() bool {
	cmd := exec.Command("vagrant", "box", "list")

	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[util/vagrant] cmd.StdoutPipe() failed", err.Error())
	}

	//
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "nanobox/boot2docker") {
				return true
			}
		}
	}()

	return false
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {

	hasBox()

	//
	if err := cmd.Start(); err != nil {
		fmt.Println("BOINK!", err)
	}

	//
	if err := cmd.Wait(); err != nil {
		fmt.Println("BOONK!", err)
	}

	if !need {
		// only on update, remove old box

		return
	}

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

	latestRelease := release{}
	latest, err := semver.NewVersion("0.0.0")
	if err != nil {
		fmt.Println("BUZUNK!", err)
	}

	//
	for _, release := range releases {
		if release.Version != "" {

			ver, err := semver.NewVersion(strings.Replace(release.Version, "v", "", -1))
			if err != nil {
				fmt.Println("BRUOONK!", err)
			}

			if latest.LessThan(*ver) {
				latestRelease = release
				latest = ver
			}
		}
	}

	fmt.Printf("LATEST RELEASE %#v\n", latestRelease)

	boxfile := filepath.Clean(config.Root + "/nanobox-boot2docker.box")

	if _, err := os.Stat(boxfile); err != nil {
		fmt.Println("LATEST RELEASE!", latestRelease)

		//
		res, err := http.Get(latestRelease.Assets[0].DownloadURL)
		if err != nil {
			config.Fatal("[commands/update] http.NewRequest() failed", err.Error())
		}
		defer res.Body.Close()

		var buf bytes.Buffer
		var percent float64
		var down int

		// format the response content length to be more 'friendly'
		total := float64(res.ContentLength) / math.Pow(1024, 2)

		// create a 'buffer' to read into
		p := make([]byte, 2048)

		//
		fmt.Printf(stylish.SubBullet("- Downloading latest CLI from %v", latestRelease.Assets[0].DownloadURL))
		for {

			// read the response body (streaming)
			n, err := res.Body.Read(p)

			// write to our buffer
			buf.Write(p[:n])

			// update the total bytes read
			down += n

			// update the percent downloaded
			percent = (float64(down) / float64(res.ContentLength)) * 100

			// show how download progress:
			// down/totalMB [*** progress *** %]
			fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

			// detect EOF and break the 'stream'
			if err != nil {
				if err == io.EOF {
					fmt.Println("")
					break
				} else {
					config.Fatal("[commands/update] res.Body.Read() failed", err.Error())
				}
			}
		}

		// replace the existing CLI with the new CLI
		fmt.Printf(stylish.SubBullet("- Replacing CLI at %s", boxfile))
		if err := ioutil.WriteFile(boxfile, buf.Bytes(), 0755); err != nil {
			if os.IsPermission(err) {
				fmt.Printf(stylish.SubBullet("[x] FAILED"))
				fmt.Printf("\nNanobox needs your permission to update.\nPlease run this command with sudo/admin privileges")
				os.Exit(0)
			}
		}

		// download the thang
		// f, err := os.Create(boxfile)
		// if err != nil {
		// 	config.Fatal("[config/update] os.Create() failed", err.Error())
		// }
		// defer f.Close()
	} else {
		fmt.Println("ALREADY HAVE!")
	}

	//
	if err := exec.Command("vagrant", "box", "add", "--name", "nanobox/boot2docker", config.Root+"/nanobox-boot2docker.box").Run(); err != nil {
		fmt.Println("BGLOINK 2!", err)
	}

	// curver, err := ioutil.ReadFile("/Users/sdomino/.nanobox/.box")
	// if err != nil {
	// 	fmt.Println("BRIZZLE!", err)
	// }
	//
	// currentVersion, err := semver.NewVersion(string(curver))
	// if err != nil {
	// 	fmt.Println("BRUNKS", err)
	// }
	//
	// fmt.Println("CURRENT VER", currentVersion, latest)

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
