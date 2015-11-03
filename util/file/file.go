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

//
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

		// skip any hidden files
		if strings.HasPrefix(fi.Name(), ".") {
			return nil
		}

		// create header for this file
		header := &tar.Header{
			Name:    file,
			Size:    fi.Size(),
			Mode:    int64(fi.Mode()),
			ModTime: fi.ModTime(),
		}

		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open the file for taring...
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy from file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})
}

//
func Untar(dest string, r io.Reader) {

	//
	gzr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	defer gzr.Close()

	//
	tr := tar.NewReader(gzr)

	//
	for {
		header, err := tr.Next()

		//
		switch {
		case err == io.EOF:
			break
		case err != nil:
			panic(err)
		}

		//
		path := filepath.Join(dest, header.Name)

		switch header.Typeflag {

		// if its a dir, make it
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				panic(err)
			}

		// if its a file, add it to the dir
		case tar.TypeReg:
			f, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			// copy from tar reader into file
			if _, err := io.Copy(f, tr); err != nil {
				panic(err)
			}

		//
		default:
			fmt.Printf("Can't: %c, %s\n", header.Typeflag, path)
		}
	}

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
		config.Fatal("[util/file/file] ioutil.ReadAll() failed - ", err.Error())
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
