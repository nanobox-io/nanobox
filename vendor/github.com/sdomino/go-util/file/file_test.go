package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	fileutil "github.com/sdomino/go-util/file"
)

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
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}
	defer tarball.Close()

	// create a tarball from src
	if err := fileutil.Tar(src, tarball); err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// varify that tarball was create at dst
	if _, err := os.Stat(filepath.Join(dst, "file.tar.gz")); err != nil {
		t.Error("Expected file, got nothing - ", err.Error())
	}

	//
	archive, err := os.Open(filepath.Join(dst, "file.tar.gz"))
	if err != nil {
		t.Error("Unexpected error - ", err.Error())
	}
	defer archive.Close()

	// untar tarball
	if err := fileutil.Untar(dst, archive); err != nil {
		t.Error("Unexpected error - ", err.Error())
	}

	// iterate through each path checking to see if file.txt was untared correctly;
	// the contents of the tarball should be the src directory (inside the dst dir)
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(dst, path, file)); err != nil {
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
