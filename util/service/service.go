//
package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// MountLocal creates a local mount (~/.nanobox/apps/<app>/<service>/<mount>)
func MountLocal() (mountName, mountDir string, err error) {

	// parse the boxfile and see if there is an service declared; if none is declared
	// simply return.
	servicePath := config.ParseBoxfile().Build.Service
	if servicePath == "" {
		return
	}

	//
	mountName = filepath.Base(servicePath)

	// if no local service is found return since there is nothing more to do here;
	// when an service is specified but not found, it's assumed that the desired
	// service exists on nanobox.io
	if _, err = os.Stat(servicePath); err != nil {
		return
	}

	//
	servicefile := filepath.Join(servicePath, "Servicefile")

	// ensure there is an servicefile at the service location
	if _, err = os.Stat(servicefile); err != nil {
		err = fmt.Errorf("No servicefile found at '%v', Exiting...\n", servicePath)
		return
	}

	// if there is an servicefile attempt to parse it to get any additional build
	// files for mounting
	files := []string{"./bin", "./Servicefile", "./meta.json"}
	if err = config.ParseConfig(servicefile, files); err != nil {
		err = fmt.Errorf("Nanobox failed to parse your Servicefile. Please ensure it is valid YAML and try again.\n")
		return
	}

	// directory to mount
	mountDir = filepath.Join(config.AppDir, mountName)

	// if the servicefile parses successfully create the mount only if it doesn't
	// already exist
	if err = os.MkdirAll(mountDir, 0755); err != nil {
		return
	}

	var abs string

	//
	abs, err = filepath.Abs(servicePath)
	if err != nil {
		return
	}

	// pull the remaining service files into the mount
	for _, files := range files {

		path := filepath.Join(abs, files)

		// just skip any files that aren't found; any required files will be
		// caught before publishing, here it doesn't matter
		if _, err := os.Stat(path); err != nil {
			continue
		}

		// copy service file into mount
		if err = fileutil.Copy(path, mountDir); err != nil {
			return
		}
	}

	return
}

// Create
func Create(name string) string {

	fmt.Printf(stylish.SubTaskStart("Creating new service on nanobox.io"))

	//
	serviceConfig := &api.ServiceConfig{
		Name: name,
	}

	//
	service, err := api.CreateService(serviceConfig)
	if err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to create service (%v).", err))
		os.Exit(1)
	}

	// wait until service has been successfuly created before continuing...
	for {
		fmt.Print(".")

		e, err := api.GetService(api.UserSlug, name)
		if err != nil {
			config.Fatal("[commands/publish] api.GetService failed", err.Error())
		}

		// once the service is "active", break
		if e.State == "active" {
			break
		}

		//
		<-time.After(1 * time.Second)
	}

	stylish.Success()

	//
	return service.ID
}

// Get gets an service from nanobox.io
func Get(userslug, name, version string) (*http.Response, error) {

	//
	if _, err := api.GetService(userslug, name); err != nil {
		os.Stderr.WriteString(stylish.ErrBullet("No official service, or service for that user found."))
		return nil, err
	}

	// if no version is provided, fetch the latest release
	if version == "" {
		os.Stderr.WriteString(stylish.ErrBullet("Please specify the version of the service you would like to fetch"))
		os.Exit(1)
	}

	//
	path := fmt.Sprintf("http://api.nanobox.io/v1/services/%v/%v/download", name, version)

	// if a user is found, pull the service from their services
	if userslug != "" {
		path = fmt.Sprintf("http://api.nanobox.io/v1/services/%v/%v/%v/download", userslug, name, version)
	}

	os.Stderr.WriteString(stylish.Bullet("Fetching service at '%s'", path))

	//
	return http.Get(path)
}

// ExtractArchive splits args on "/" looking for a user and archive:
// - user/service-name
// - user/service-name=0.0.1
func ExtractArchive(s string) (user, archive string) {

	split := strings.Split(s, "/")

	// switch on the length to determine if the split resulted in a user and a service
	// or just an service
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
		// fmt.Printf("%v is not a valid format when fetching an service (see help).\n", args[0])
		os.Exit(1)
	}

	return
}

// ExtractService splits on the archive to find the service and the release (version)
func ExtractService(archive string) (service, version string) {

	// split on '=' looking for a version
	split := strings.Split(archive, "=")

	// switch on the length to determine if the split resulted in a service and version
	// or just an service
	switch len(split) {

	// if len is 1 then just an service was found (no version specified)
	case 1:
		service = split[0]

	// if len is 2 then an service and version were found
	case 2:
		service = split[0]
		version = split[1]
	}

	return
}
