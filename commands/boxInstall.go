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
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var boxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: boxInstall,
}

func hasBox() bool {

	need := true

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
				need = false
			}
		}
	}()

	//
	if err := cmd.Start(); err != nil {
		fmt.Println("BOINK!", err)
	}

	//
	if err := cmd.Wait(); err != nil {
		fmt.Println("BOONK!", err)
	}

	return need
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {

	hasBox()

	boxmd5path := "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/boot2docker.box.md5"
	boxfilemd5 := filepath.Clean(config.Root + "/nanobox-boot2docker.box.md5")

	boxpath := "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/boot2docker.box"
	boxfile := filepath.Clean(config.Root + "/nanobox-boot2docker.box")

	if _, err := os.Stat(boxfilemd5); err != nil {

		//
		res, err := http.Get(boxpath)
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
		fmt.Printf(stylish.SubBullet("- Downloading latest CLI from %v", boxpath))
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

		//
		//
		//
		//
		md5res, err := http.Get(boxmd5path)
		if err != nil {
			config.Fatal("[commands/update] http.NewRequest() failed", err.Error())
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(md5res.Body)
		if err != nil {
			fmt.Println("BOIOIOIOIOINK", err)
		}

		if err := ioutil.WriteFile(boxfilemd5, b, 0755); err != nil {
			if os.IsPermission(err) {
				fmt.Printf(stylish.SubBullet("[x] FAILED"))
				fmt.Printf("\nNanobox needs your permission to update.\nPlease run this command with sudo/admin privileges")
				os.Exit(0)
			}
		}

	} else {
		fmt.Println("ALREADY HAVE!")
	}

	//
	if err := exec.Command("vagrant", "box", "add", "--name", "nanobox/boot2docker", boxfile).Run(); err != nil {
		fmt.Println("BGLOINK 2!", err)
	}
}
