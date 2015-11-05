// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package file

import (
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
	src, dst := setup()
	defer cleanup(src, dst)

	file := "file.txt"
	paths := []string{
		"",
		"deep",
		"deep/deep",
		"deep/deep/deep",
	}

	//
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Join(src, path), 0755); err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(filepath.Join(src, path, file), []byte("contents"), 0644); err != nil {
			panic(err)
		}
	}

	// copy file from src to dst
	Copy(dst, src)

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
	src, dst := setup()
	defer cleanup(src, dst)

	file := "file.txt"
	paths := []string{
		"",
		"deep",
		"deep/deep",
		"deep/deep/deep",
	}

	//
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Join(src, path), 0755); err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(filepath.Join(src, path, file), []byte("contents"), 0644); err != nil {
			panic(err)
		}
	}

	//
	tarball, err := os.Create(filepath.Join(dst, "file.tar.gz"))
	defer tarball.Close()
	if err != nil {
		panic(err)
	}

	// create a tarball from src
	if err := Tar(src, tarball); err != nil {
		panic(err)
	}

	// varify that tarball was create at dst
	if _, err := os.Stat(filepath.Join(dst, "file.tar.gz")); err != nil {
		t.Error("Expected file, got nothing - ", err.Error())
	}

	//
	archive, err := os.Open(filepath.Join(dst, "file.tar.gz"))
	defer archive.Close()
	if err != nil {
		panic(err)
	}

	// untar tarball
	if err := Untar(dst, archive); err != nil {
		panic(err)
	}

	// iterate through each path checking to see if file.txt was untared correctly;
	// the contents of the tarball should be the src directory
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(src, path, file)); err != nil {
			t.Error("Expected file, got nothing - ", err.Error())
		}
	}
}

//
func newDir(name string) (string, error) {
	return ioutil.TempDir("", name)
}

//
func cleanup(files ...string) {
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			panic(err)
		}
	}
}

//
func setup() (src, dst string) {

	var err error

	//
	if src, err = newDir("src"); err != nil {
		panic(err)
	}

	//
	if dst, err = newDir("dst"); err != nil {
		panic(err)
	}

	return
}
