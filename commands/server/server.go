package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/jcelliott/lumber"

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
		log.Fatal("server needs to run as privilaged user")
		return
	}
	// make sure things know im the server
	registry.Set("server", true)

	// set the logger on linux and osx to go to /var/log
	if runtime.GOOS != "windows" {
		fileLogger, err := lumber.NewTruncateLogger("/var/log/nanobox.log")
		if err != nil {
			log.Printf("logging error:%s\n", err)
		}

		//
		lumber.SetLogger(fileLogger)
	}

	// fire up the service manager (only required on windows)
	go svcStart()

	go startEcho()

	go updateUpdater()

	// first up the tap driver (only required on osx)
	go startTAP()

	// add any registered rpc classes
	for _, controller := range registeredRPCs {
		rpc.Register(controller)
	}

	// only listen for rpc calls on localhost
	listener, e := net.Listen("tcp", "127.0.0.1:23456")
	if e != nil {
		log.Fatal("listen error:", e)
		return
	}

	// listen for new connections forever
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go rpc.ServeConn(conn)
		}
	}
}

// TEMP: this only ever needs to be run once.
// but it wont hurt to run it once everytime nanobox server starts
// this can be removed once everyone is >= 2.1.0
func updateUpdater() {
	// update.Server = true
	update.Name = strings.Replace(update.Name, "nanobox", "nanobox-update", 1)
	update.TmpName = strings.Replace(update.TmpName, "nanobox", "nanobox-update", 1)
	path, err := exec.LookPath(os.Args[0])
	path = strings.Replace(path, "nanobox", "nanobox-update", 1)
	if err != nil {
		return
	}
	log.Println(update.Run(path))
}

// run a client request to the rpc server
func ClientRun(funcName string, args interface{}, response interface{}) error {
	// log.Printf("clientcall: %s %#v %#v\n", funcName, args, response)
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
func startTAP() {
	if runtime.GOOS == "darwin" {
		exec.Command("/sbin/kextload", "/Library/Extensions/tap.kext").Run()
	}
}

type handle struct {
}

func (handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong\n"))
}

func startEcho() {
	http.ListenAndServe(":65000", handle{})
}
