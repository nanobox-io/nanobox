package helpers

import "github.com/nanobox-io/nanobox/models"

func Endpoint(envModel *models.Env, args []string) ([]string, string, string) {
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

	return args, "production", "default"
}
