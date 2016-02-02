//
package service

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/auth"
	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	s3util "github.com/nanobox-io/nanobox/util/s3"
	serviceutil "github.com/nanobox-io/nanobox/util/service"
)

//
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes a service to nanobox.io",
	Long:  ``,

	Run: publish,
}

// publish
func publish(ccmd *cobra.Command, args []string) {
	stylish.Header("publishing service")

	// ensure there is an ./Servicefile
	if _, err := os.Stat("./Servicefile"); err != nil {
		fmt.Printf("No Servicefile found. Be sure to publish from a project directory. Exiting...\n")
		os.Exit(1)
	}

	// create a new service
	fmt.Printf(stylish.Bullet("Creating release..."))
	serviceConfig := &api.ServiceConfig{}

	// parse the ./Servicefile into the new service
	if err := config.ParseConfig("./Servicefile", serviceConfig); err != nil {
		fmt.Printf("Nanobox failed to parse your Servicefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	fmt.Printf(stylish.Bullet("Verifying service is publishable..."))

	// determine if any required fields (name, version, language, summary) are missing,
	// if any are found to be missing exit 1
	// NOTE: I do this using fallthrough for asthetics only. The message is generic
	// enough that all cases will return the same message, and this looks better than
	// a single giant case/if (var == "" || var == "" || ...)
	switch {
	case serviceConfig.Category == "":
		fallthrough
	case serviceConfig.Image == "":
		fallthrough
	case serviceConfig.Name == "":
		fallthrough
	case serviceConfig.Summary == "":
		fallthrough
	case serviceConfig.Version == "":
		fmt.Printf(stylish.Error("required fields missing", `Your Servicefile is missing one or more of the following required fields for publishing:

  name:      # the name of your project
  summary:   # a 140 character summary of the project
	category:  #
	image:     #

Please ensure all required fields are provided and try again.`))

		os.Exit(1)
	}

	// attempt to read a README.md file and add it to the release...
	b, err := ioutil.ReadFile("./README.md")
	if err != nil {

		// this only fails if the file is not found, EOF is not an error. If no Readme
		// is found exit 1
		fmt.Printf(stylish.Error("missing readme", "Your service is missing a README.md file. This file is required for publishing, as it is the only way for you to communicate how to use your service. Please add a README.md and try again."))
		os.Exit(1)
	}

	//
	serviceConfig.Readme = string(b)

	// check to see if the service already exists on nanobox.io
	fmt.Printf(stylish.Bullet("Checking for existing service on nanobox.io"))

	//
	api.UserSlug, api.AuthToken = auth.Authenticate()

	// if no service exists create a new one
	service, err := api.GetService(api.UserSlug, serviceConfig.Name)
	if err != nil {
		if apiErr, _ := err.(api.APIError); apiErr.Code == 404 {
			service.ID = serviceutil.Create(serviceConfig.Name)
		} else {
			config.Fatal("[commands/service/publish] api.GetService failed", err.Error())
		}
	}

	// create a meta.json file where we can add any extra data we might need; since
	// this is only used for internal purposes the file is removed once we're done
	// with it
	meta, err := os.Create("./meta.json")
	if err != nil {
		config.Fatal("[commands/service/publish] os.Create() failed", err.Error())
	}
	defer meta.Close()
	defer os.Remove(meta.Name())

	// add any custom info to the metafile
	meta.WriteString(fmt.Sprintf(`{"service_id": "%s"}`, service.ID))

	// list of required files/folders for an service
	files := []string{"./bin", "./Servicefile", "./meta.json"}

	// check to ensure no required files are missing from build folder
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			fmt.Printf(stylish.Error("required files missing", fmt.Sprintf("Unable to find %s; this file is either required, or was declared as a required file in the Servicefile. Please read the following documentation to ensure all required files are included and try again. \n\ndocs.nanobox.io/services/project-creation/#example-service-file-structure\n", file)))
			os.Exit(1)
		}
	}

	// create the temp services folder for building the tarball
	tarPath := filepath.Join(config.ServicesDir, serviceConfig.Name)
	if err := os.MkdirAll(tarPath, 0755); err != nil {
		config.Fatal("[commands/service/publish] os.Create() failed", err.Error())
	}

	// remove tarDir once published
	defer func() {
		if err := os.RemoveAll(tarPath); err != nil {
			os.Stderr.WriteString(stylish.ErrBullet("Failed to remove '%v'...", tarPath))
		}
	}()

	// parse the ./Servicefile again to get all build files
	if err := config.ParseConfig("./Servicefile", files); err != nil {
		fmt.Printf("Nanobox failed to parse your Servicefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	// add each of the build files to the final tarPath; not handling the error here
	// because it simply means the file doesn't exist and therefor wont be copied
	// to the final tarball
	for _, file := range files {
		fileutil.Copy(file, tarPath)
	}

	// create an empty buffer for writing the file contents to for the subsequent
	// upload
	archive := bytes.NewBuffer(nil)

	//
	h := md5.New()

	// create a tarball to upload the represents the service
	if err := fileutil.Tar(tarPath, archive, h); err != nil {
		config.Fatal("[commands/service/publish] file.Tar() failed", err.Error())
	}

	// create a checksum for the new release once its finished being archived
	serviceConfig.Checksum = fmt.Sprintf("%x", h.Sum(nil))

	//
	// attempt to upload the release to S3
	fmt.Printf(stylish.Bullet("Uploading release to s3..."))

	v := url.Values{}
	v.Add("user_slug", api.UserSlug)
	v.Add("auth_token", api.AuthToken)
	v.Add("version", serviceConfig.Version)

	//
	s3url, err := s3util.RequestURL(fmt.Sprintf("http://api.nanobox.io/v1/services/%v/request_upload?%v", serviceConfig.Name, v.Encode()))
	if err != nil {
		config.Fatal("[commands/service/publish] s3.RequestURL() failed", err.Error())
	}

	// upload the release to s3
	if err := s3util.Upload(s3url, archive); err != nil {
		config.Fatal("[commands/service/publish] s3.Upload() failed", err.Error())
	}

	//
	// if the release uploaded successfully to s3, created one on odin
	fmt.Printf(stylish.Bullet("Uploading release to nanobox.io"))
	if _, err := api.CreateService(serviceConfig); err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to publish release (%v).", err))
		os.Exit(1)
	}
}
