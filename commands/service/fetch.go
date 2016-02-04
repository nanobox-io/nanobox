//
package service

import (
	"fmt"
	"io"
	"os"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	// "github.com/nanobox-io/nanobox/auth"
	"github.com/nanobox-io/nanobox/config"
	serviceutil "github.com/nanobox-io/nanobox/util/service"
)

//
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches an service from nanobox.io",
	Long: `
Description:
  Fetches an service from nanobox.io

  Allowed formats when fetching an service
  - service-name
  - service-name=0.0.1
  - user/service-name
  - user/service-name=0.0.1
	`,

	Run: fetch,
}

//
func init() {

	// no default is set here because we define the value later, once we know the
	// name and version of the service they are fetching
	fetchCmd.Flags().StringVarP(&fFile, "ouput-document", "O", "", "specify where to save the service")
}

// fetch
func fetch(ccmd *cobra.Command, args []string) {

	//
	// api.UserSlug, api.AuthToken = auth.Authenticate()

	if len(args) == 0 {
		os.Stderr.WriteString("Please provide the name of an service you would like to fetch, (run 'nanobox service fetch -h' for details)")
		os.Exit(1)
	}

	os.Stderr.WriteString(stylish.Bullet("Attempting to fetch '%v'", args[0]))

	// extract a user and archive (desired service) from args[0]
	user, archive := serviceutil.ExtractArchive(args[0])

	// extract a service and version from the archive
	service, version := serviceutil.ExtractService(archive)

	// pull the service from nanobox.io
	res, err := serviceutil.Get(user, service, version)
	if err != nil {
		config.Fatal("[commands/service/fetch] http.Get() failed", err.Error())
	}
	defer res.Body.Close()

	//
	switch res.StatusCode / 100 {
	case 2, 3:
		break
	case 4:
		os.Stderr.WriteString(stylish.ErrBullet("No service at version '%v' found", version))
		os.Exit(1)
	case 5:
		os.Stderr.WriteString(stylish.ErrBullet("Failed to fetch service (%v).", res.Status))
		os.Exit(1)
	}

	// determine if destination will be a file or stdout (stdout by default)
	dest := os.Stdout
	defer dest.Close()

	// write the download to the local file system
	if fFile != "" {

		//
		f, err := os.Create(fFile)
		if err != nil {
			os.Stderr.WriteString(stylish.ErrBullet("Unable to save file - %v", err.Error()))
			os.Stderr.WriteString("Exiting...\n")
			return
		}

		// if the file was created successfully then set it as the destination
		os.Stderr.WriteString(stylish.Bullet("Saving service as '%s'", fFile))
		dest = f
	}

	// write the file
	if _, err := io.Copy(dest, res.Body); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("[commands.fetch] io.Copy() failed%s", err.Error()))
	}
}
