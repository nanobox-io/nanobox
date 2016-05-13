package commands

import (
  "fmt"
  "os"
  "io"

  "github.com/spf13/cobra"
  "github.com/nanobox-io/nanobox-golang-stylish"
  "github.com/nanobox-io/golang-docker-client"

  "github.com/nanobox-io/nanobox/provider"
  "github.com/nanobox-io/nanobox/util/print"
)

var (

  TestCmd = &cobra.Command{
    Use: "test",
    Short: "Just testing the output possibilities",
    Long: ``,
    Run: func(ccmd *cobra.Command, args []string) {
      testOutput()
    },
  }
)

func testOutput() {
  fmt.Print(stylish.NestedBullet("Provisioning Platform Services...", 0))
  fmt.Print(stylish.NestedBullet("Launching load balancer", 1))
  fmt.Print(stylish.NestedProcessStart("Downloading docker image nanobox/portal", 2))

  provider.DockerEnv()
  docker.Initialize("env")

  pr, pw := io.Pipe()
  go print.DisplayJSONMessagesStream(pr, os.Stdout, os.Stdout.Fd(), true, stylish.GenerateNestedPrefix(3), nil)
  docker.ImagePull("alpine:latest", pw)
  fmt.Print(stylish.ProcessEnd())

  fmt.Print(stylish.NestedBullet("Starting docker container...", 2))
}

//
// + Booting nanobox vm --------------------------->
//   output
//   output
//   blablablal
//   blablabal
//
// + Provisioning Platform Services...
//   - Launching load balancer
//   | - Downloading docker image nanobox/portal -->
//   | | output
//   | | output
//   | | output
//   | | blablablabla
//   |
//   - Launching warehouse...
