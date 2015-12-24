//
package vagrant

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nanobox-io/nanobox/config"
)

//
var err error

// Exists ensure vagrant is installed
func Exists() (exists bool) {
	var err error

	// check if vagrant is installed
	if _, err = exec.LookPath("vagrant"); err == nil {

		// initilize Vagrant incase it hasn't been; there is a chance that Vagrant has
		// never been used meaning there won't be a .vagrant.d folder, so we initialize
		// vagrant just to ensure it's ready to be used with nanobox by running any
		// vagrant command (in this case "vagrant -v").
		if b, err := exec.Command("vagrant", "-v").CombinedOutput(); err != nil {
			config.Fatal("[util/vagrant/vagrant] exec.Command() failed", string(b))
		}

		// read setup_version to determine if the version of vagrant is too old
		// (< 1.5.0) and needs to be migrated
		b, err := ioutil.ReadFile(filepath.Join(config.Home, ".vagrant.d", "setup_version"))
		if err != nil {
			config.Fatal("[util/vagrant/vagrant] ioutil.ReadFile() failed", err.Error())
		}

		// convert the []byte value from the file into a float 'version'
		version, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			config.Fatal("[util/vagrant/vagrant] strconv.ParseFloat() failed", err.Error())
		}

		// if the current version of vagrant is less than a 'working version' (1.5)
		// give instructions on how to update
		if version < 1.5 {
			fmt.Println(`
Nanobox has detected that you are using an old version of Vagrant (<1.5). Before
you can continue you'll need to run "vagrant update" and follow the instructions
to update Vagrant.
			`)

			// exit here to allow for upgrade
			os.Exit(0)
		}

		// if all checks pass
		exists = true
	}

	return
}

// run runs a vagrant command
func run(cmd *exec.Cmd) error {

	//
	handleCMDout(cmd)

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output above
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// runInContext runs a command in the context of a Vagrantfile (from the same dir)
func runInContext(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<config.App>. Running the command from
	// the directory that contains the Vagratfile ensure that the command can
	// atleast run (especially in cases like 'create' where a VM hadn't been created
	// yet, and a UUID isn't available)
	setContext(config.AppDir)

	//
	handleCMDout(cmd)

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output above
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	// switch back to project dir
	setContext(config.CWDir)

	return nil
}

// setContext changes the working directory to the designated context
func setContext(context string) {
	if err := os.Chdir(context); err != nil {
		fmt.Printf("No app found at %s. Exiting...\n", config.AppDir)
		os.Exit(1)
	}
}

func customScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {

	//
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}

	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// handleCMDout
func handleCMDout(cmd *exec.Cmd) {

	// create a stderr pipe that will write any error messages to the log
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Fatal("[util/vagrant/vagrant] cmd.StderrPipe() failed", err.Error())
	}

	// log any command errors to the log
	stderrScanner := bufio.NewScanner(stderr)
	go func() {

		//
		var once sync.Once
		for stderrScanner.Scan() {

			// only display the error message once, but log every error
			once.Do(func() { Error("A vagrant error occured", "") })
			Log.Error(stderrScanner.Text())
		}
	}()

	// create a stdout pipe that will allow for scanning the output line-by-line;
	// if needed a stderr pipe could also be created at some point
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Fatal("[util/vagrant/vagrant] cmd.StdoutPipe() failed", err.Error())
	}

	// scan the command output intercepting only 'important' lines of vagrant output'
	// and tailoring their message so as to not flood the output.
	// styled according to: http://nanodocs.gopagoda.io/engines/style-guide
	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(customScanner)

	// this is a mapping of all the vagrant output and our output to show instead
	filter := map[string]string{
		"VM not created. Moving on...":                   "Nanobox not yet created, use 'nanobox dev' or 'nanobox run' to create it.",
		"VirtualBox VM is already running.":              "This nanobox is already running",
		"Importing base box 'nanobox/boot2docker'...":    "Importing nanobox base image",
		"Booting VM...":                                  "Booting virtual machine",
		"Configuring and enabling network interfaces...": "Configuring virtual network",
		"Mounting shared folders...":                     fmt.Sprintf("Mounting source code (%s)", config.CWDir),
		"Waiting for nanobox server...":                  "Starting nanobox server",
		"Attempting graceful shutdown of VM...":          "Shutting down virtual machine",
		"Forcing shutdown of VM...":                      "Shutting down virtual machine",
		"Saving VM state and suspending execution...":    "Saving virtual machine",
		// "Resuming suspended VM...":                       "Resuming virtual machine",   // this is handled by commands/resume.go
		// "Destroying VM and associated drives...":         "Destroying virtual machine", // this is handled by commands/destroy.go
	}

	var done chan bool
	var next chan string
	progressing := false

	//
	go func() {
		for stdoutScanner.Scan() {

			//
			txt := strings.TrimSpace(stdoutScanner.Text())

			// log all vagrant output (might as well)
			Log.Info(txt)

			//
			switch {

			// show the progress bar if vagrant downloads nanobox/boot2docker
			case strings.Contains(txt, "box: Progress:"):

				progressing = true

				// subMatch[1] - percentage downloaded
				// subMatch[2] - amount/s
				// subMatch[3] - estimated time remaining
				subMatch := regexp.MustCompile(`box: Progress: (\d{1,3})% \(Rate: (.*), Estimated time remaining: (\d*:\d*:\d*)`).FindStringSubmatch(txt)

				// if for some reason we don't get the matches we need just skip that
				// line; this should never happen unless vagrant change their ouput
				if len(subMatch) < 4 {
					continue
				}

				// if for some reason the submatch fails to convert to an int just skip
				// the line; this will likely never happen
				i, err := strconv.Atoi(subMatch[1])
				if err != nil {
					continue
				}

				// show download progress: [*** progress *** 0.0%] 00:00:00 remaining
				fmt.Printf("\r\033[K   [%-41s %s%%] %s (%s remaining)", strings.Repeat("*", int(float64(i)/2.5)), subMatch[1], subMatch[2], subMatch[3])

			// if it's NOT a progress string and we were previously progressing
			case !strings.Contains(txt, "box: Progress:") && progressing:
				fmt.Printf("\r\033[K   [**************************************** 100.0%%] 00:00:00 remaining\n")
				progressing = false

				// fallthrough so we dont miss any messages
				fallthrough

			// filter and print any messages received from vagrant
			case strings.HasPrefix(txt, fmt.Sprintf("==> %v: ", config.Nanofile.Name)):
				subMatch := regexp.MustCompile(`==> \S*:\s(.*)`).FindStringSubmatch(txt)

				// if for some reason the regex failed to pull anything just skip the
				// line; this should never happen unless vagrant changes their output
				if len(subMatch) <= 1 {
					continue
				}

				// if the line is found in our filter print it
				if v, ok := filter[subMatch[1]]; ok {

					// indicate that there the next message has arrived by closing the previous
					// message channel
					if next != nil {
						close(next)
					}

					// wait on done (skipping the first wait)
					if done != nil {
						<-done
					}

					// make a new next
					next = make(chan string)

					// make a new done
					done = make(chan bool)

					// fire up the printer/dotter
					go printMessage(next, done)

					// send the next message
					next <- v
				}
			}
		}

		// close next once the scan is complete
		if next != nil {
			close(next)
		}

		// close done once the scan is complete
		// if done != nil {
		// 	close(done)
		// }
	}()
}

// printMessage takes a message and done channel, reading of the message channel
// and printing the message with a loader
func printMessage(msg chan string, done chan bool) {

	var out string

	for {
		select {

		// read the message
		case msg, ok := <-msg:

			// if !ok that means the channel was closed (another message received) so
			// print the final output, close the spinner and done, and return
			if !ok {
				fmt.Printf("\r\033[K   - %s\n", out)

				// stop the spinner
				done <- true

				close(done)
				return
			}

			// because a final message ("") will come when the channel is closed, store
			// the message here after that happens to have the actual messages
			out = msg

			// fire up the spinner
			go func() {

				spinner := `-\|/`
				i := 0

				for {
					select {

					// spin baby spin
					default:
						fmt.Printf("\r\033[K   %s %s", string(spinner[i%len(spinner)]), out)
						<-time.After(time.Second / 24)
						i++

					// message complete, stop spinner
					case <-done:
						return
					}
				}
			}()
		}
	}
}
