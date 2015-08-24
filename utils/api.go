// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package utils

import (
	"errors"
	"fmt"
	"regexp"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/ui"
)

// HandleAPIError takes an error returned from an API call, break it down and
// return important information regarding the error. The Pagoda Box API returns
// custom errors in some instances that need to have very specific handlers.
func HandleAPIError(err error) error {

	// if its a pagodabox.Error we have special things we want to do...
	if apiError, ok := err.(api.APIError); ok {

		//
		switch apiError.Code {

		// Unauthorized, Forbidden, Not Found, Internal Server Error, Bad Gateway
		case 401, 403, 404, 500, 502:
			return errors.New(apiError.Body)

		// Unprocessable Entity -
		case 422:

			// separate the custom 422 error from the message (ex. {"upgrade":["Cannot exceed free limit"]})
			reFindError := regexp.MustCompile(`^\{\s*\"(.*)\"\s*\:\s*\[\s*\"(.*)\"\s*\]\s*\}$`)

			subMatch := reFindError.FindStringSubmatch(apiError.Body)
			if subMatch == nil {
				ui.LogFatal("[utils/api] reFindError.FindStringSubmatch() failed", errors.New(fmt.Sprintf("No matches found for api error: %v", apiError.Body)))
			}

			return errors.New(fmt.Sprintf("[utils/api] %d %v - %v", 422, subMatch[1], subMatch[2]))

		// some error we're not aware of
		default:
			return errors.New(fmt.Sprintf("[utils/api] Unhandled API error - %v", err))
		}

		// ...if not, just write to the log
	} else {
		return errors.New(fmt.Sprintf("[utils/api] Unhandled error - %v", err))
	}

	return nil
}
