// Package dns ...
package dns

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nanobox-io/nanobox/commands/server"
)

type DomainName struct {
	IP     string
	Domain string
}

type DomainRPC struct{}

type Request struct {
	Entry string
}

type Response struct {
	Message string
	Success bool
}

var (
	hostsFile = detectHostsFile()
	newline   = detectNewlineChar()
)

func init() {
	server.Register(&DomainRPC{})
}

// Entry generate the DNS entry to be added
func Entry(ip, name, env string) string {
	return fmt.Sprintf("%s     %s # dns added for '%s' by nanobox", ip, name, env)
}

// Exists ...
func Exists(entry string) bool {

	// open the hosts file for scanning...
	f, err := os.Open(hostsFile)
	if err != nil {
		return false
	}
	defer f.Close()

	// scan each line of the hosts file to see if there is a match for this
	// entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == entry {
			return true
		}
	}

	return false
}

func List(filter string) []DomainName {

	// open the hosts file
	f, err := os.Open(hostsFile)
	if err != nil {
		return nil
	}
	defer f.Close()

	entries := []DomainName{}

	// scan each line of the hosts file to see if there is a match for this
	// entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), filter) {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 2 {
				entries = append(entries, DomainName{IP: fields[0], Domain: fields[1]})
			}
		}
	}

	return entries
}

// Add ...
func Add(entry string) error {
	// break early if there is no entry
	// or we have already added this entry
	if entry == "" || Exists(entry) {
		return nil
	}

	req := Request{entry}
	resp := &Response{}

	err := server.ClientRun("DomainRPC.Add", req, resp)
	if !resp.Success {
		err = fmt.Errorf("failed to add domain: %v %v", err, resp.Message)
	}

	return err
}

// the rpc function run from the server
func (drpc *DomainRPC) Add(req Request, resp *Response) error {

	// open hosts file
	f, err := os.OpenFile(hostsFile, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the DNS entry to the file
	// we place a newline before and after because
	// the extra new lines wont hurt and it ensures success
	if _, err := f.WriteString(fmt.Sprintf("%s%s%s", newline, req.Entry, newline)); err != nil {
		return err
	}

	resp.Success = true
	return nil
}

// Remove ...
func Remove(entry string) error {
	if entry == "" {
		return nil
	}
	req := Request{entry}
	resp := &Response{}

	err := server.ClientRun("DomainRPC.Remove", req, resp)
	if !resp.Success {
		err = fmt.Errorf("failed to remove domain: %v %v", err, resp.Message)
	}
	return nil
}

// the rpc function run from the server
func (drpc *DomainRPC) Remove(req Request, resp *Response) error {

	// "contents" will end up storing the entire contents of the file excluding the
	// entry that is trying to be removed
	var contents string

	// open hosts file
	f, err := os.OpenFile(hostsFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// remove entry from /etc/hosts
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if the line contain the entry skip it
		// make it do a loose string check
		// if its exactly the entry then remove it.
		// if it contains the same environment
		// also remove it
		if strings.Contains(scanner.Text(), req.Entry) {
			continue
		}

		// add each line back into the file
		contents += fmt.Sprintf("%s%s", scanner.Text(), newline)
	}

	// trim the contents to avoid any extra newlines
	contents = strings.TrimSpace(contents)

	// add a single newline for completeness
	contents += newline

	// write back the contents of the hosts file minus the removed entry
	if err := ioutil.WriteFile(hostsFile, []byte(contents), 0644); err != nil {
		return err
	}

	resp.Success = true
	return nil

}

// RemoveAll removes all dns entries added by nanobox
func RemoveAll() error {

	// short-circuit if no entries were added by nanobox
	if len(List("by nanobox")) == 0 {
		return nil
	}

	return Remove("by nanobox")
}
