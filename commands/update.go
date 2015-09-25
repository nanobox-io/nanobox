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

	"github.com/pagodabox/nanobox-cli/util"
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

	//
	path, err := osext.Executable()
	if err != nil {
		util.Fatal("[commands/update] osext.ExecutableFolder() failed", err)
	}

	// download a new CLI from s3 that matches their os and arch
	download := fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%v/%v/nanobox", runtime.GOOS, runtime.GOARCH)

	// create a new request
	fmt.Printf(stylish.SubBullet("- Downloading latest CLI from %v", download))
	req, err := http.NewRequest("GET", download, nil)
	if err != nil {
		util.Fatal("[commands/update] http.NewRequest() failed", err)
	}

	// download the new CLI
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.Fatal("[commands/update] http.DefaultClient.Do() failed", err)
	}

	var buf bytes.Buffer
	p := make([]byte, 2048)

	var percent float64

	down := 0
	total := float64(res.ContentLength) / math.Pow(1024, 2)

	// replace the existing CLI with the new CLI
	for {

		// stream read the download
		n, err := res.Body.Read(p)

		// write to our buffer
		buf.Write(p[:n])

		// update the total bytes read
		down += n

		percent = (float64(down) / float64(res.ContentLength)) * 100

		// show how download progress
		fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

		if err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			} else {
				util.Fatal("[commands/update] res.Body.Read() failed", err)
			}
		}

		defer res.Body.Close()
	}

	//
	fmt.Printf(stylish.SubBullet("- Replacing CLI at %s", path))
	ioutil.WriteFile(path, buf.Bytes(), 0755)

	//
	fmt.Println(stylish.Success())

	//
	// config.Console.Debug("[commands.update] command: %v, args: %+v", os.Args[0], os.Args[1:])

	// // attempt to run the command that was being run to begin with (unless its update)
	// if os.Args[1] != "update" {

	//  // run the command that was being run to begin with
	//  out, err := exec.Command(os.Args[0], os.Args[1:]...).Output()
	//  if err != nil {
	//    util.Fatal("[commands.update] exec.Command()", err)
	//  }

	//  // show the output of the command that is run
	//  fmt.Printf("%v\n", out)
	// }
}
