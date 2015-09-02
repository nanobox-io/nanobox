// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"archive/tar"
	// "compress/gzip"
	"io"
	"os"
	// "path/filepath"
	// "sync"

	// "github.com/pagodabox/nanobox-cli/auth"
	// "github.com/pagodabox/nanobox-cli/ui"
)

var tw *tar.Writer

//
func Gzip() {

}

//
func Tar() {

}

//
func CreateTarBall() {

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
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open the file for taring...
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy the file data to the tarball
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}

	return nil
}
