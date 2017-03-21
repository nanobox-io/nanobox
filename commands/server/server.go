package server

import (
	"net"
	"log"
	"net/rpc"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a dedicated nanobox server",
	Long:  ``,
	Run: serverFnc,
}

func serverFnc(ccmd *cobra.Command, args []string) {
	if !util.IsPrivileged() {
		log.Fatal("server needs to run as privilaged user")	
		return
	}
	// start the rpc server
	rpc.Register(commands)
	listener, e := net.Listen("tcp", ":23456")
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
