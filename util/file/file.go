// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package file

import (
	"archive/tar"
	// "bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	// "path/filepath"
	// "sync"
	"math"
	"strings"

	// "github.com/nanobox-io/nanobox-cli/auth"
	// "github.com/nanobox-io/nanobox-cli/ui"
	// "github.com/nanobox-io/nanobox-cli/config"
)

var GZ *gzip.Writer
var TW *tar.Writer

//
func Gzip() {

}

//
func Tar() {

}

//
func TarBall() {

}

// tarFile
func tarFile(path string, fi os.FileInfo, err error) error {

	// only want to tar files...
	if !fi.Mode().IsDir() {

		// fmt.Println("TARING!", path)

		// create header for this file
		header := &tar.Header{
			Name:    path,
			Size:    fi.Size(),
			Mode:    int64(fi.Mode()),
			ModTime: fi.ModTime(),
		}

		// write the header to the tarball archive
		if err := TW.WriteHeader(header); err != nil {
			return err
		}

		// open the file for taring...
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy the file data to the tarball
		if _, err := io.Copy(TW, f); err != nil {
			return err
		}
	}

	return nil
}

// Download
func Download(path string, w io.Writer) error {
	res, err := http.Get(path)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

// Progress
func Progress(path string, w io.Writer) error {

	//
	download, err := http.Get(path)
	if err != nil {
		return err
	}
	defer download.Body.Close()

	var percent float64
	var down int

	// format the response content length to be more 'friendly'
	total := float64(download.ContentLength) / math.Pow(1024, 2)

	// create a 'buffer' to read into
	p := make([]byte, 2048)

	//
	for {

		// read the response body (streaming)
		n, err := download.Body.Read(p)

		// write to our buffer
		w.Write(p[:n])

		// update the total bytes read
		down += n

		// update the percent downloaded
		percent = (float64(down) / float64(download.ContentLength)) * 100

		// show download progress: down/totalMB [*** progress *** %]
		fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

		// detect EOF and break the 'stream'
		if err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			} else {
				return err
			}
		}
	}

	return nil
}
