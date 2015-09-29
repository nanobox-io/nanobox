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

//
func needBox() bool {

	need := true

	cmd := exec.Command("vagrant", "box", "list")

	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[commands/boxInstall] cmd.StdoutPipe() failed", err.Error())
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
		config.Fatal("[commands/boxInstall] cmd.Start() failed", err.Error())
	}

	//
	if err := cmd.Wait(); err != nil {
		config.Fatal("[commands/boxInstall] cmd.Wait() failed", err.Error())
	}

	return need
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {

	// check to see if a box even needs to be installed
	if getMD5() == compairMD5() && !needBox() {
		return
	}

	boxpath := "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.box"
	boxfile := filepath.Clean(config.Root + "/nanobox-boot2docker.box")
	if _, err := os.Stat(boxfile); err != nil {

		//
		box, err := http.Get(boxpath)
		if err != nil {
			config.Fatal("[commands/boxInstall] http.Get() failed", err.Error())
		}
		defer box.Body.Close()

		var buf bytes.Buffer
		var percent float64
		var down int

		// format the response content length to be more 'friendly'
		total := float64(box.ContentLength) / math.Pow(1024, 2)

		// create a 'buffer' to read into
		p := make([]byte, 2048)

		//
		fmt.Printf(stylish.SubBullet("- Downloading virtual machine image from %v", boxpath))
		for {

			// read the response body (streaming)
			n, err := box.Body.Read(p)

			// write to our buffer
			buf.Write(p[:n])

			// update the total bytes read
			down += n

			// update the percent downloaded
			percent = (float64(down) / float64(box.ContentLength)) * 100

			// show how download progress:
			// down/totalMB [*** progress *** %]
			fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

			// detect EOF and break the 'stream'
			if err != nil {
				if err == io.EOF {
					fmt.Println("")
					break
				} else {
					config.Fatal("[commands/boxInstall] res.Body.Read() failed", err.Error())
				}
			}
		}

		//
		if err := ioutil.WriteFile(boxfile, buf.Bytes(), 0755); err != nil {
			config.Fatal("[commands/boxInstall] ioutil.WriteFile() failed", err.Error())
		}

		//
		md5, err := http.Get("https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5")
		if err != nil {
			config.Fatal("[commands/boxInstall] http.Get() failed", err.Error())
		}
		defer md5.Body.Close()

		b, err := ioutil.ReadAll(md5.Body)
		if err != nil {
			config.Fatal("[commands/boxInstall] ioutil.ReadAll() failed", err.Error())
		}

		if err := ioutil.WriteFile(config.Root+"/nanobox-boot2docker.md5", b, 0755); err != nil {
			config.Fatal("[commands/boxInstall] ioutil.WriteFile() failed", err.Error())
		}
	}

	//
	if err := exec.Command("vagrant", "box", "add", "--force", "--name", "nanobox/boot2docker", boxfile).Run(); err != nil {
		config.Fatal("[commands/boxInstall] exec.Command() failed", err.Error())
	}
}
