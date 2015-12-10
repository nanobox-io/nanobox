//
package engine

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
)

// RemountLocal simply calls MountLocal() but only returns the error as the
// mounts name and dir are not important in this instance.
func RemountLocal() (err error) {
	_, _, err = MountLocal()
	return
}

// MountLocal creates a local mount (~/.nanobox/apps/<app>/<engine>/<mount>)
func MountLocal() (mountName, mountDir string, err error) {

	// parse the boxfile and see if there is an engine declared; if none is declared
	// simply return.
	enginePath := config.ParseBoxfile().Build.Engine
	if enginePath == "" {
		return
	}

	//
	mountName = filepath.Base(enginePath)

	// if no local engine is found return since there is nothing more to do here;
	// when an engine is specified but not found, it's assumed that the desired
	// engine exists on nanobox.io
	if _, err = os.Stat(enginePath); err != nil {
		return
	}

	//
	enginefile := filepath.Join(enginePath, "Enginefile")

	// ensure there is an enginefile at the engine location
	if _, err = os.Stat(enginefile); err != nil {
		err = fmt.Errorf("No enginefile found at '%v', Exiting...\n", enginePath)
		return
	}

	// if there is an enginefile attempt to parse it
	mount := &struct {
		Overlays []string `json:"overlays"`
	}{}

	if err = config.ParseConfig(enginefile, mount); err != nil {
		err = fmt.Errorf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		return
	}

	mountDir = filepath.Join(config.AppDir, mountName)

	// if the enginefile parses successfully create the mount only if it doesn't
	// already exist
	if _, err = os.Stat(mountDir); err != nil {
		if err = os.Mkdir(mountDir, 0755); err != nil {
			return
		}
	}

	// iterate through each overlay and untar it to the mount dir
	for _, overlay := range mount.Overlays {
		GetOverlay(overlay, mountDir)
	}

	var abs string

	abs, err = filepath.Abs(enginePath)
	if err != nil {
		return
	}

	// pull the remainin engine files into the mount
	for _, f := range []string{"bin", "Enginefile", "lib", "templates", "files"} {

		path := filepath.Join(abs, f)

		// just skip any files that aren't found; any required files will be
		// caught before publishing, here it doesn't matter
		if _, err := os.Stat(path); err != nil {
			continue
		}

		// copy engine file into mount
		if err = fileutil.Copy(path, mountDir); err != nil {
			return
		}
	}

	return
}

// GetEngine gets an engine from nanobox.io
func GetEngine(user, archive, version string) (*http.Response, error) {

	//
	engine, err := api.GetEngine(user, archive)
	if err != nil {
		os.Stderr.WriteString(stylish.ErrBullet("No official engine, or engine for that user found."))
		return nil, err
	}

	// if no version is provided, fetch the latest release
	if version == "" {
		version = engine.ActiveReleaseID
	}

	//
	path := fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/releases/%v/download", archive, version)

	// if a user is found, pull the engine from their engines
	if user != "" {
		path = fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/%v/releases/%v/download", user, archive, version)
	}

	os.Stderr.WriteString(stylish.Bullet("Fetching engine at '%s'", path))

	//
	return http.Get(path)
}

// GetOverlay fetches an overlay from nanobox.io and untars it into dst
func GetOverlay(overlay, dst string) {

	// extract a user and archive (desired engine)
	user, archive := ExtractArchive(overlay)

	// extract an engine and version from the archive
	engine, version := ExtractEngine(archive)

	//
	res, err := GetEngine(user, engine, version)
	if err != nil {
		config.Fatal("[util/engine/engine] http.Get() failed", err.Error())
	}
	defer res.Body.Close()

	//
	switch res.StatusCode / 100 {
	case 2, 3:
		break
	case 4, 5:
		os.Stderr.WriteString(stylish.ErrBullet("Unable to fetch '%v' overlay. Exiting...\n", engine))
		os.Exit(1)
	}

	//
	if err := fileutil.Untar(dst, res.Body); err != nil {
		config.Fatal("[util/engine/engine] file.Untar() failed", err.Error())
	}
}

// ExtractArchive splits args on "/" looking for a user and archive:
// - user/engine-name
// - user/engine-name=0.0.1
func ExtractArchive(s string) (user, archive string) {

	split := strings.Split(s, "/")

	// switch on the length to determine if the split resulted in a user and a engine
	// or just an engine
	switch len(split) {

	// if len is 1 then only a download was found (no user specified)
	case 1:
		archive = split[0]

		// if len is 2 then a user was found (from which to pull the download)
	case 2:
		user = split[0]
		archive = split[1]

	// any other number or args
	default:
		// fmt.Printf("%v is not a valid format when fetching an engine (see help).\n", args[0])
		os.Exit(1)
	}

	return
}

// ExtractEngine splits on the archive to find the engine and the release (version)
func ExtractEngine(archive string) (engine, version string) {

	// split on '=' looking for a version
	split := strings.Split(archive, "=")

	// switch on the length to determine if the split resulted in a engine and version
	// or just an engine
	switch len(split) {

	// if len is 1 then just an engine was found (no version specified)
	case 1:
		engine = split[0]

	// if len is 2 then an engine and version were found
	case 2:
		engine = split[0]
		version = split[1]
	}

	return
}
