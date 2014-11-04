package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/vaughan0/go-ini"

	"github.com/nanobox-core/cli/ui"
)

// FindGitConfig attempts to find a .git/config file and return the file and
// containing path. If no file is found (path '/'), returns an error
func FindGitConfigFile() (ini.File, string, error) {

	// get the current working directory (cwd) to use when looking for a .git dir
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get cwd. See ~/.pagodabox/log.txt for details")
		ui.Error("helpers.FindGitConfigFile", err)
	}

	// attempt to find a .git dir, starting from cwd
	path := findGitDir(cwd)
	if path == "/" {
		return nil, cwd, errors.New("Unable to detect .git folder. It appears you are not in a git repository...")
	}

	// attempt to load the config file
	gitConfig, err := ini.LoadFile(path + "/.git/config")
	if err != nil {
		fmt.Println("Unable to load .git/config file. See ~/.pagodabox/log.txt for details")
		ui.Error("helpers.FindGitConfigFile", err)
	}

	return gitConfig, cwd, nil
}

// findGitDir attempts to reverse recursively find a .git directory starting from
// 'path'
func findGitDir(path string) string {
	f, _ := os.Stat(path + "/.git")

	// end of the line and we haven't found anything.
	if path == "/" {
		return path
	}

	// If no file is found go up a dir and look again.
	if f == nil {
		return findGitDir(filepath.Dir(path))
	}

	// if a '.git' folder is found the parent folder path is returned, otherwise
	// '/' is returned
	return path
}

// FindPagodaApp attempts to find a .git/config file and iterate over the 'remotes'
// looking for any 'pagoda' remotes. If none are found, prompts for an app-name.
// If one is found, it returns that app-name. If more than one is found, display
// a list of found apps and prompt for which to use.
func FindPagodaApp() string {

	reMatchRemote := regexp.MustCompile("git@git.pagodabox.io:apps/")
	reFindRemote := regexp.MustCompile(`^git@git\.pagodabox\.io\:apps\/(.*)\.git$`)

	apps := []string{}
	remotes := 0

	// attempt to find a git config file. Don't need to worry about 'path' or 'err'
	// here, only interested in the file
	gitConfigFile, _, _ := FindGitConfigFile()

	// count how many Pagoda Box remotes there are (if any)
	for name, _ := range gitConfigFile {

		section, ok := gitConfigFile.Get(name, "url")

		if ok && reMatchRemote.MatchString(section) {
			remotes++

			subMatch := reFindRemote.FindStringSubmatch(section)
			if subMatch == nil {
				fmt.Println("Unable to parse remote. See ~/.pagodabox/log.txt for details")
				ui.Error("helper:FindPagodaApp", errors.New("No matches found for remote: "+section))
			}

			apps = append(apps, subMatch[1])
		}
	}

	switch {

	// no Pagoda Box remotes found, prompt for an app to use
	case remotes <= 0:
		fmt.Println("We were unable to find a Pagoda Box app tied to this project.\n")
		return ui.Prompt("Which app are you trying to use: ")

	// one Pagoda Box remote found, return the app-name
	case remotes == 1:
		return apps[0]

	// multiple Pagoda Box remotes found, display apps found and prompt for an app
	// to use
	case remotes > 1:
		fmt.Println("We found the following Pagoda Box apps tied to this project.\n")
		for _, app := range apps {
			fmt.Println("- " + app)
		}
		fmt.Println("")
		return ui.Prompt("Which app would you like to use: ")
	}

	return ""
}
