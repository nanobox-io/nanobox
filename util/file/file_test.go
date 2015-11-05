// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package file

import (
	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// tests if Copy works as intended; a deep file is tested with the assumption that
// if it can be successfully copied over, than anything else along the way would
// have also
func TestCopy(t *testing.T) {

	// create tmp dirs
	src, dst, file, paths, err := setup()
	defer cleanup(src, dst)
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// copy file from src to dst
	Copy(src, dst)

	// iterate through each path checking to see if file.txt was copied
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(dst, path, file)); err != nil {
			t.Error("Expected file, got nothing - ", err.Error())
		}
	}

}

// tests if Tar and Untar are working as intended; can't really test taring w/o
// also testing untaring
func TestTarUntar(t *testing.T) {

	// create tmp dirs
	src, dst, file, paths, err := setup()
	defer cleanup(src, dst)
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	//
	tarball, err := os.Create(filepath.Join(dst, "file.tar.gz"))
	defer tarball.Close()
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// create a tarball from src
	if err := Tar(src, tarball); err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// varify that tarball was create at dst
	if _, err := os.Stat(filepath.Join(dst, "file.tar.gz")); err != nil {
		t.Error("Expected file, got nothing - ", err.Error())
	}

	//
	archive, err := os.Open(filepath.Join(dst, "file.tar.gz"))
	defer archive.Close()
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// untar tarball
	if err := Untar(dst, archive); err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// iterate through each path checking to see if file.txt was untared correctly;
	// the contents of the tarball should be the src directory (inside the dst dir)
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(dst, src, path, file)); err != nil {
			t.Error("Expected file, got nothing - ", err.Error())
		}
	}
}

// setup creates a source and destination directory and fills the source with some
// files/folders
func setup() (src, dst, file string, paths []string, err error) {

	file = "file.txt"
	paths = []string{
		"",
		"deep",
		"deep/deep",
		"deep/deep/deep",
	}

	//
	if src, err = ioutil.TempDir("", "src"); err != nil {
		return src, dst, file, paths, err
	}

	//
	if dst, err = ioutil.TempDir("", "dst"); err != nil {
		return src, dst, file, paths, err
	}

	//
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Join(src, path), 0755); err != nil {
			return src, dst, file, paths, err
		}

		if err := ioutil.WriteFile(filepath.Join(src, path, file), []byte("contents"), 0644); err != nil {
			return src, dst, file, paths, err
		}
	}

	return
}

// cleanup removes the src and dst dirs
func cleanup(files ...string) (err error) {
	for _, file := range files {
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}

	return
}
