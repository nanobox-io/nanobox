package helpers

import "github.com/nanobox-io/nanobox/models"

func Endpoint(envModel *models.Env, args []string, maxArgs int) ([]string, string, string) {
	if len(args) == 0 {
		return args, "production", "default"
	}

	switch args[0] {
	case "local":
		return args[1:], "local", "dev"
	case "dry-run":
		return args[1:], "local", "sim"
	default:
		_, ok := envModel.Remotes[args[0]]
		if ok {
			return args[1:], "production", args[0]
		}
	}

	// if we were given the maximum number of arguments then the first artument must be a production
	// application name that was not in our remotes
	if maxArgs == len(args) {
		return args[1:], "production", args[0]
	}

	// todo: MAYBE check the remote here (`fetch the remote` in other locations in code)
	// // fetch the remote
	// remote, ok := envModel.Remotes[args[0]]
	// if ok {
	// 	// set the app id
	//  return args[1:], "production", remote.ID
	// }

	return args, "production", "default"
}
