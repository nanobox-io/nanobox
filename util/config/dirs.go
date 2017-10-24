package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-boxfile"
)

// GlobalDir ...
func GlobalDir() string {

	// set Home based off the users homedir (~)
	p, err := homedir.Dir()
	if err != nil {
		return ""
	}

	globalDir := filepath.ToSlash(filepath.Join(p, ".nanobox"))
	os.MkdirAll(globalDir, 0755)

	return globalDir
}

// LocalDir returns the current working directory
func LocalDir() string {
	var cwd string
	var err error

	// use syscall if possible (os uses $PWD, causing irregularities on darwin)
	if syscall.ImplementsGetwd {
		cwd, err = syscall.Getwd()
		if err != nil {
			fmt.Printf("Failed to syscall.getwd - %s\n", err.Error())

			// fallback to os.Getwd
			cwd, err = os.Getwd()
			if err != nil {
				fmt.Printf("Failed to os.getwd - %s\n", err.Error())
				return ""
			} else {
				return cwd
			}
			return ""
		}
	} else {
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Printf("Failed to os.getwd - %s\n", err.Error())
			// todo: this func should return an error and setup should fail if this fails
			return ""
		}
	}

	// return filepath.ToSlash(cwd)

	// fmt.Println("CURRENT DIR: ", cwd)

	// todo: boxfile validation should have happened previously. need to find out why it doesn't
	boxfilePresent := func(path string) bool {
		boxfile := filepath.ToSlash(filepath.Join(path, "boxfile.yml"))
		fi, err := os.Stat(boxfile)
		if err != nil {
			return false
		}
		return !fi.IsDir()
	}

	path := cwd
	for !boxfilePresent(path) {
		if path == "" || path == "/" || strings.HasSuffix(path, ":\\") {
			// return the current working directory if we cant find a path
			// todo: this returns the path, which isn't what we want
			return filepath.ToSlash(cwd)
		}
		// eliminate the most child directory and then check it
		path = filepath.Dir(path)
	}

	// recursively check for boxfile

	return filepath.ToSlash(path)
}

// LocalDirName ...
func LocalDirName() string {
	return filepath.Base(LocalDir())
}

// SSHDir ...
func SSHDir() string {

	//
	p, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return filepath.ToSlash(filepath.Join(p, ".ssh"))
}

var engineDir string

// EngineDir gets the directory of the engine if it is a directory and on the
// local file system
func EngineDir() (string, error) {
	if engineDir != "" {
		return engineDir, nil
	}

	engineName := boxfile.NewFromPath(Boxfile()).Node("run.config").StringValue("engine")
	if engineName != "" {
		// if engine specified is not a filepath, set engine path to blank
		var validLocal = regexp.MustCompile(`^[~|\.|\/|\\]`)
		if !validLocal.MatchString(engineName) {
			return "", nil
		}

		fi, err := os.Stat(engineName)
		if err != nil {
			return "", fmt.Errorf("Failed to find engine - %s", err.Error())
		}
		if fi.IsDir() {
			path, err := filepath.Abs(engineName)
			if err != nil {
				return "", fmt.Errorf("Failed to resolve engine location - %s", err.Error())
			}
			engineDir = path
			return path, nil
		}
	}

	return "", nil
}

// BinDir creates a directory where nanobox specific binaries can be downloaded
// docker, dockermachine, etc
func BinDir() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		if runtime.GOOS == "windows" {
			return "C:\\Program Files\\Nanobox"

		}
		return filepath.Dir(os.Args[0])
	}
	return filepath.Dir(path)
}

func EtcDir() string {

	etcDir := filepath.ToSlash(filepath.Join(GlobalDir(), "etc"))

	if err := os.MkdirAll(etcDir, 0755); err != nil {
		lumber.Fatal("[config/config] os.Mkdir() failed", err.Error())
	}

	return etcDir
}
