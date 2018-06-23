package server

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/update"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a dedicated nanobox server",
	Long:  ``,
	Run:   serverFnc,
}

const name = "nanobox-server"

func serverFnc(ccmd *cobra.Command, args []string) {
	if !util.IsPrivileged() {
		fmt.Println("server needs to run as a privileged user")
		return
	}
	// make sure things know im the server
	registry.Set("server", true)

	// set the logger on linux and osx to go to /var/log
	if runtime.GOOS != "windows" {
		fileLogger, err := lumber.NewTruncateLogger("/var/log/nanobox.log")
		if err != nil {
			fmt.Printf("logging error:%s\n", err)
		}

		lumber.SetLogger(fileLogger)
	}

	lumber.Info("Starting nanobox server...")

	// fire up the service manager (only required on windows)
	go svcStart()

	go startEcho()

	go updateUpdater()

	// first up the tap driver (only required on osx)
	lumber.Info("Attempting to load tap driver...")
	err := startTAP()
	if err != nil {
		lumber.Info("Failed to load tap driver - %s", err.Error())
		// err 102 for nanobox means kext failed to load
		os.Exit(102)
	}
	lumber.Info("Tap driver loaded.")

	// add any registered rpc classes
	for _, controller := range registeredRPCs {
		rpc.Register(controller)
	}

	lumber.Info("Attempting to listen on port 23456...")
	// only listen for rpc calls on localhost
	listener, err := net.Listen("tcp", "127.0.0.1:23456")
	if err != nil {
		lumber.Info("Failed to listen - %s", err.Error())
		return
	}

	lumber.Info("Nanobox server listening")

	// listen for new connections forever
	for {
		if conn, err := listener.Accept(); err != nil {
			lumber.Fatal("accept error: " + err.Error())
		} else {
			lumber.Info("new connection established\n")
			go rpc.ServeConn(conn)
		}
	}
}

// updateUpdater used to be a temporary means to update everyone's updater,
// but it is quite useful so we will leave it in. Maybe in the future we'll
// try updating nanobox itself prior to starting.
func updateUpdater() {
	lumber.Info("Updating nanobox-update")
	update.Name = strings.Replace(update.Name, "nanobox", "nanobox-update", 1)
	update.TmpName = strings.Replace(update.TmpName, "nanobox", "nanobox-update", 1)

	// this gets the path to nanobox (assumes nanobox-update is at same location)
	lumber.Info("Attempting to find nanobox - %s", os.Args[0])
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		lumber.Info("Failed to find nanobox - %s", err.Error())
		return
	}
	path = strings.Replace(path, "nanobox", "nanobox-update", 1)
	lumber.Info("Updating - %s", path)
	err = update.Run(path)
	if err != nil {
		lumber.Info("Failed to update `nanobox-update` - %s", err.Error())
		return
	}
	lumber.Info("Update complete")
}

// run a client request to the rpc server
func ClientRun(funcName string, args interface{}, response interface{}) error {
	// lumber.Info("clientcall: %s %#v %#v\n", funcName, args, response)
	client, err := rpc.Dial("tcp", "127.0.0.1:23456")
	if err != nil {
		return err
	}

	err = client.Call(funcName, args, response)
	if err != nil {
		return err
	}

	return nil
}

// the tap driver needs to be loaded anytime nanobox is running the vpn (always on osx)
func startTAP() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	out, err := exec.Command("kextstat").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to show loaded kernel extensions - %s. Output - %s", err.Error(), out)
	}

	if !strings.Contains(string(out), "net.sf.tuntaposx.tap") {
		lumber.Info("Loading tap extension...")
		out, err := exec.Command("kextload", "/Library/Extensions/tap.kext").CombinedOutput()
		if err != nil {
			return fmt.Errorf("Failed to load tap extensions - %s. Output - %s", err.Error(), out)
		}
		lumber.Info("Loaded tap extension.")
	}

	return nil
}

type handle struct {
}

func (handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong\n"))
}

func startEcho() {
	http.ListenAndServe(":65000", handle{})
}
