// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import (
	"bytes"
	"flag"
	"github.com/gorilla/pat"
	"github.com/nanobox-io/nanobox/config"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

var child = flag.Bool("child", false, "")

func run(name string, handle func()) *exec.Cmd {
	if *child {
		handle()
		os.Exit(0)
	} else {
		return exec.Command("go", "test", "-run", name, ".", "-child")
	}
	return nil
}

func startServer(test *testing.T, handler http.Handler) io.Closer {
	test.Log("listening")
	listen, err := net.Listen("tcp", config.ServerURI)
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	test.Log("starting to serve")
	go http.Serve(listen, handler)
	test.Log("served")

	return listen
}

func TestExec(test *testing.T) {
	config.ServerURI = "127.0.0.1:1234"
	child := run("TestExec", func() {
		in := bytes.NewReader([]byte("nothing"))
		out := bytes.NewBuffer([]byte{})
		err := execInternal("command", "cmd=ls", in, out)
		test.Log(err)
	})

	mux := pat.New()
	mux.Get("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("pong"))
	})
	mux.Post("/exec", func(res http.ResponseWriter, req *http.Request) {
		conn, _, err := res.(http.Hijacker).Hijack()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		defer conn.Close()
		script := req.FormValue("cmd")
		if script == "" {
			test.Log("missing script")
			test.FailNow()
		}

		cmd := exec.Command("/bin/bash", "-c", script)
		cmd.Stdout = io.MultiWriter(conn, os.Stdout)
		cmd.Stdin = conn
		cmd.Stderr = io.MultiWriter(conn, os.Stderr)
		err = cmd.Run()
		conn.Close()
	})
	listen := startServer(test, mux)
	defer listen.Close()
	errChan := make(chan error)
	go func() {
		output, err := child.CombinedOutput()
		test.Log(string(output))
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()
	select {
	case <-time.After(time.Second * 4):
		test.Log("timed out...")
		test.FailNow()
	case err := <-errChan:
		if err == nil {
			return
		}
		test.Log("child failed to run")
		test.Log(err)
		test.FailNow()
	}
}
