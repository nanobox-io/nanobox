package helpers

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/ui"
)

// HandleAPIError takes an error returned from an API call, break it down and
// return important information regarding the error. The Pagoda Box API returns
// custom errors in some instances that need to have very specific handlers.
func HandleAPIError(err error) (int, string, string) {

	// if its a pagodabox.Error we have special things we want to do...
	if apiError, ok := err.(pagodaAPI.Error); ok {

		switch apiError.Code {

		// These should be rare. so we'll just print the error and exit
		case 401, 403, 500, 502:
			fmt.Printf(apiError.Status)
			os.Exit(1)

		// Not Found - Resouce does not exist
		case 404:
			return 404, apiError.StatusText, "Reason: does not exist"

		// Unprocessable Entity -
		case 422:

			// separate the custom 422 error from the message (ex. {"upgrade":["Cannot exceed free limit"]})
			reFindError := regexp.MustCompile(`^\{\s*\"(.*)\"\s*\:\s*\[\s*\"(.*)\"\s*\]\s*\}$`)

			subMatch := reFindError.FindStringSubmatch(apiError.Body)
			if subMatch == nil {
				fmt.Println("Unable to parse api error. See ~/.pagodabox/log.txt for details")
				ui.Error("error:HandleAPIError", errors.New("No matches found for api error: "+apiError.Body))
			}

			return 422, subMatch[1], subMatch[2]

		// some error we're not aware of
		default:
			fmt.Println("Unhandled API error. See ~/.pagodabox/log.txt for details")
			ui.Error("helpers.HandleAPIError", err)
		}

		// ...if not, just panic
	} else {
		fmt.Println("Unhandled error. See ~/.pagodabox/log.txt for details")
		ui.Error("helpers.HandleAPIError", err)
	}

	return 1, "", ""
}
