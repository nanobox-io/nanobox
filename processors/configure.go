package processors

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

var configured bool

func Configure() error {
	var os string

	// make sure to only run configure one time
	if configured {
		return nil
	}
	configured = true

	v, err := util.OsDetect()
	if err == nil {
		os = v
	}

	if os == "high sierra" && !models.HasRead() {
		// warn about high sierra
		hasRead := stringAsker(`
--------------------------------------------------------------------------------
+ WARNING:
+
+ MacOS High Sierra introduces breaking changes to Nanobox!
+
+ Please ensure you have read the following guides before continuing:
+ https://content.nanobox.io/installing-nanobox-on-macos-high-sierra/
--------------------------------------------------------------------------------

Have you already read the guide? y/n`, map[string]string{"y": "yes", "n": "no"})
		if hasRead == "no" {
			exec.Command("open", "https://content.nanobox.io/installing-nanobox-on-macos-high-sierra/").Start()
			return fmt.Errorf("\nEnding configure, please read the guide and try again.\n")
		}
		models.DoneRead()
	}

	// todo: why do we wait?
	<-time.After(150 * time.Millisecond)

	config := &models.Config{
		Provider:  "docker-machine",
		MountType: "native",
		CPUs:      1,
		RAM:       1,
	}

	fmt.Print(`
CONFIGURE NANOBOX
---------------------------------------------------------------
Please answer the following questions so we can customize your
nanobox configuration. Feel free to update your config at any
time by running: 'nanobox configure'

(Learn more at : https://docs.nanobox.io/local-config/configure-nanobox/)
`)

	defer func() {
		fmt.Println(`
      **
   *********
***************   [√] Nanobox successfully Configured!
:: ********* ::   ------------------------------------------------------------
" ::: *** ::: "   Change these settings at any time via : 'nanobox configure'
  ""  :::  ""
    "" " ""
       "
`)
	}()

	// ask about provider
	config.Provider = stringAsker(`
How would you like to run nanobox?
  a) Inside a lightweight VM
  b) Via Docker Native

  Note : Mac users, we strongly recommend choosing (a) until Docker Native
         resolves an issue causing slow speeds : http://bit.ly/2jYFfWQ

Answer: `, map[string]string{"a": "docker-machine", "b": "native"})

	// if provider == docker-machine ask more questions
	if config.Provider == "native" {
		config.Save()
		return nil
	}

	// ask about cpus
	config.CPUs = intAsker(fmt.Sprintf(`
How many CPU cores would you like to make available to the VM (1-%d)?
-------------------------------------------------------------------
  Note : we recommend 2 or more

Answer: `, runtime.NumCPU()), runtime.NumCPU())

	// ask about ram
	config.RAM = intAsker(`
How many GB of RAM would you like to make available to the VM (2-4)?
-------------------------------------------------------------------
  Note : we recommended 2 or more

Answer: `, 8)

	if os != "high sierra" {
		// ask about mount types
		config.MountType = stringAsker(`
Would you like to enable netfs for faster filesystem access (y/n)?
-------------------------------------------------------------------
  Note : We HIGHLY recommend (y). Using this option may prompt for password

Answer: `, map[string]string{"y": "netfs", "n": "native"})
	}

	config.Save()

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
