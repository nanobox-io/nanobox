// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"runtime"
	"strings"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the CLI to the newest available version",
	Long: `
Description:
  Updates the CLI to the newest available version`,

	Run: nanoUpdate,
}

// nanoUpdate
func nanoUpdate(ccmd *cobra.Command, args []string) {

	fmt.Printf(stylish.Bullet("Updating nanobox CLI"))

	// download a new CLI from s3 that matches their os and arch
	download := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH)

	//
	res, err := http.Get(download)
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
	fmt.Printf(stylish.SubBullet("- Downloading latest CLI from %v", download))
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

	// get the path of the current executing CLI
	path, err := osext.Executable()
	if err != nil {
		config.Fatal("[commands/update] osext.ExecutableFolder() failed", err.Error())
	}

	// replace the existing CLI with the new CLI
	fmt.Printf(stylish.SubBullet("- Replacing CLI at %s", path))
	if err := ioutil.WriteFile(path, buf.Bytes(), 0755); err != nil {
		fmt.Println("ERR! %#v", err)
	}

	//
	fmt.Printf(stylish.SubBullet("- Now running %s", config.VERSION))

	//
	fmt.Println(stylish.Success())

	// // attempt to run the command that was being run to begin with (unless its update)
	// if os.Args[1] != "update" {

	//  // run the command that was being run to begin with
	//  out, err := exec.Command(os.Args[0], os.Args[1:]...).Output()
	//  if err != nil {
	//    config.Fatal("[commands.update] exec.Command()", err.Error())
	//  }

	//  // show the output of the command that is run
	//  fmt.Printf("%v\n", out)
	// }
}
