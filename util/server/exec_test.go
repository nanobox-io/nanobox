// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/kr/pty"
	"github.com/nanobox-io/nanobox/config"
	"io"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func startServer(test *testing.T, handler http.Handler) io.Closer {
	listen, err := net.Listen("tcp", config.ServerURI)
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	go http.Serve(listen, handler)

	return listen
}

func normalPing(mux *pat.Router) {
	mux.Get("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("pong"))
	})
}

func normalExec(test *testing.T, mux *pat.Router) {
	mux.Post("/exec", func(res http.ResponseWriter, req *http.Request) {
		test.Log("got exec")
		conn, rw, err := res.(http.Hijacker).Hijack()
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
		test.Log("executing", script)
		cmd := exec.Command("/bin/bash", "-c", script)
		cmd.Stdin = rw
		cmd.Stderr = conn
		// cmd.Stdout = conn
		// cmd.Run()
		reader, err := cmd.StdoutPipe()
		cmd.Start()
		io.Copy(conn, reader)
		conn.(*net.TCPConn).CloseWrite()
		err = cmd.Wait()
		test.Log("finished running")
	})
}

func ptyExec(test *testing.T, mux *pat.Router) {
	mux.Post("/exec", func(res http.ResponseWriter, req *http.Request) {
		test.Log("got exec")
		conn, rw, err := res.(http.Hijacker).Hijack()
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
		test.Log("executing", script)
		cmd := exec.Command("/bin/bash", "-c", script)
		file, err := pty.Start(cmd)
		if err != nil {
			test.Log(err)
			test.FailNow()
		}
		test.Log("copying from pty")
		go io.Copy(file, rw)
		io.Copy(conn, file)
		test.Log("finished running")
	})
}

func TestExec(test *testing.T) {
	config.ServerURI = "127.0.0.1:1234"

	mux := pat.New()
	normalPing(mux)
	normalExec(test, mux)
	listen := startServer(test, mux)
	defer listen.Close()

	errChan := make(chan error)
	go func() {
		in := bytes.NewBuffer([]byte("this is a test"))
		out := &bytes.Buffer{}
		err := execInternal("exec", "command", "cmd=cat", in, out)
		if err != nil {
			errChan <- err
			return
		}
		if out.String() != "this is a test" {
			test.Log("output:", out.Len())
			errChan <- fmt.Errorf("unexpected output: '%v'", out.String())
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
		test.Log(err)
		test.FailNow()
	}
}

func TestPTYExec(test *testing.T) {
	config.ServerURI = "127.0.0.1:1234"

	mux := pat.New()
	normalPing(mux)
	ptyExec(test, mux)
	listen := startServer(test, mux)
	defer listen.Close()

	errChan := make(chan error)
	go func() {
		in := bytes.NewBuffer([]byte("this is a test\n"))
		in.Write([]byte{4}) // EOT
		out := &bytes.Buffer{}
		err := execInternal("exec", "command", "cmd=cat", in, out)
		if err != nil {
			errChan <- err
			return
		}
		// a tty does local echo, and has some extra stuff
		if regexp.MustCompile("this is a test.+this is a test.+").Match(out.Bytes()) {
			test.Log("output:", out.Len())
			errChan <- fmt.Errorf("unexpected output: '%v'", out)
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
		test.Log(err)
		test.FailNow()
	}
}

func TestExecEarlyExit(test *testing.T) {
	config.ServerURI = "127.0.0.1:1235"

	mux := pat.New()
	normalPing(mux)
	normalExec(test, mux)
	listen := startServer(test, mux)
	defer listen.Close()

	errChan := make(chan error)
	go func() {
		// need to use a pipe so that no EOF is returned. this was causing test to fail very quickly
		r, _ := io.Pipe()
		out := &bytes.Buffer{}
		err := execInternal("exec", "command", "cmd=exit", r, out)
		test.Log("exited", err)
		if err != nil {
			errChan <- err
			return
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
		test.Log(err)
		test.FailNow()
	}
}
