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
	"compress/gzip"
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Tar
func Tar(path string, writers ...io.Writer) error {

	//
	mw := io.MultiWriter(writers...)

	//
	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	//
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	//
	return filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {

		//
		if err != nil {
			return err
		}

		// only tar files (not dirs)
		if fi.Mode().IsRegular() {

			// create header for this file
			header := &tar.Header{
				Name: strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator)),
				Mode: int64(fi.Mode()),
				Size: fi.Size(),
				// ModTime:  fi.ModTime(),
				// Typeflag: tar.TypeReg,
			}

			// write the header to the tarball archive
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// open the file for taring...
			f, err := os.Open(file)
			defer f.Close()
			if err != nil {
				return err
			}

			// copy from file data into tar writer
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil
	})
}

// Untar
func Untar(dst string, r io.Reader) error {

	//
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		return err
	}

	//
	tr := tar.NewReader(gzr)

	//
	for {
		header, err := tr.Next()

		//
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		dir := filepath.Dir(header.Name)
		base := filepath.Base(header.Name)
		dirPath := filepath.Join(dst, dir)

		// if the dir doesn't exist it needs to be created
		if _, err := os.Stat(dirPath); err != nil {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return err
			}
		}

		// create the file
		f, err := os.OpenFile(filepath.Join(dirPath, base), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer f.Close()

		// copy over contents
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}
	}
}

// Download
func Download(path string, w io.Writer) error {
	res, err := http.Get(path)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		config.Fatal("[util/file/file] ioutil.ReadAll() failed - ", err.Error())
	}

	w.Write(b)

	return nil
}

// Progress
func Progress(path string, w io.Writer) error {

	//
	download, err := http.Get(path)
	defer download.Body.Close()
	if err != nil {
		return err
	}

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
