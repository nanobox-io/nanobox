package server

import (
	"log"
	"net"
	"net/rpc"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util"
)

type Response struct {
	Output   string
	ExitCode int
}

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

	// fire up the service manager (only required on windows)
	go svcStart()

	// register controllers
	rpc.Register(bridge)
	rpc.Register(commands)

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

// run a client request to the rpc server
func run(funcName string, args interface{}, response interface{}) error {
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
