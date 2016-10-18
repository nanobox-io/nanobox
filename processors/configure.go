package processors

import (
	"time"
	"os"
	"fmt"
	"runtime"

	"github.com/nanobox-io/nanobox/util/config"
)

func Configure() error {

	<-time.After(time.Second)
	
	setupConf := &config.SetupConf{
		Provider: "docker-machine",
		Mount:    "native",
		CPUs:     1,
		RAM:      1,
	}

	fmt.Print(`
CONFIGURE NANOBOX 
---------------------------------------------------------------
Please answer the following questions so we can customize your
nanobox configuration. Feel free to update your config at any 
time by running: 'nanobox configure'

(Learn more at : http://docs.nanobox.io/configure)
`)


	// ask about provider
	// currently ignoring the input here
	stringAsker(`
How would you like to run nanobox?
  a) Inside a lightweight VM
  b) Via Docker Native (coming)

Answer: `, map[string]string{"a": "docker-machine", "b": "native"})

	// if provider == docker-machine ask more questions
	if setupConf.Provider == "native" {
		config.ConfigFile(setupConf)
		return nil
	}

	// ask about cpus
	setupConf.CPUs = intAsker(fmt.Sprintf(`
How many CPU cores would you like to make available to the VM (1-%d)?

Answer: `,runtime.NumCPU()), runtime.NumCPU())

	// ask about ram
	setupConf.RAM = intAsker(`How many GB of RAM would you like to make available to the VM (2-4)?

Answer: `, 4)

	// ask about mount types
	setupConf.Mount = stringAsker(`
Would you like to enable netfs for faster filesystem access (y/n)?
(we highly recommend using this option, but this will prompt for password)

Answer: `, map[string]string{"y": "netfs", "n": "native"})


	config.ConfigFile(setupConf)
	return nil

}

func stringAsker(text string, answers map[string]string) string {
	var answer string

	fmt.Fprint(os.Stdout, text)
	fmt.Scanln(&answer)

	result, ok := answers[answer]
	for !ok {
		fmt.Println("Invalid response, please try again:")
		fmt.Print(text)
		fmt.Scanln(&answer)
		result, ok = answers[answer]
	}
	return result
}

func intAsker(text string, max int) int {
	var answer int

	fmt.Print(text)
	fmt.Scanln(&answer)

	for answer > max {
		fmt.Println("\nInvalid response, please try again:\n")
		fmt.Print(text)
		fmt.Scanln(&answer)
	}
	return answer
}
