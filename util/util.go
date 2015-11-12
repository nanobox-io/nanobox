//
package util

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

// VboxExists ensure virtualbox is installed; if ever there is a virtualbox package
// this can be moved there
func VboxExists() (exists bool) {
	if err := exec.Command("vboxmanage", "-v").Run(); err == nil {
		exists = true
	}

	return
}

// MD5sMatch determines if a local MD5 matches a remote MD5
func MD5sMatch(localFile, remotePath string) (bool, error) {

	// read the local file; will return os.PathError if doesn't exist
	b, err := ioutil.ReadFile(localFile)
	if err != nil {
		return false, err
	}

	// get local md5 checksum (as a string)
	localMD5 := fmt.Sprintf("%x", md5.Sum(b))

	// GET remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// read the remote md5 checksum
	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// compare checksum's
	return strings.TrimSpace(localMD5) == strings.TrimSpace(string(remoteMD5)), nil
}

// StringToIP generates an IPv4 address based off the app name for use as a
// vagrant private_network IP.
func StringToIP(s string) string {

	var network uint32 = 2886729728 // 172.16.0.0 network
	var sum uint32 = 0              // the last two octets of the assigned network

	// create an md5 of the app name to ensure a uniqe IP is generated each time
	h := md5.New()
	io.WriteString(h, s)

	// iterate through each byte in the md5 hash summing along the way
	for _, v := range []byte(h.Sum(nil)) {
		sum += uint32(v)
	}

	ip := make(net.IP, 4)

	// convert app name into a unique private network IP by adding the first portion
	// of the network with the generated portion
	binary.BigEndian.PutUint32(ip, (network + sum))

	return ip.String()
}

// StringToPort generates a unique network port to allow running multiple vms at
// once
func StringToPort(s string) string {

	port := 10000 // starting port is > than 100000 to try and avoid confilcts

	// create an md5 of the app name to ensure a uniqe port is generated each time
	h := md5.New()
	io.WriteString(h, s)

	// iterate through each byte in the md5 hash summing along the way
	for _, v := range []byte(h.Sum(nil)) {
		port += int(v)
	}

	return fmt.Sprint(port)
}
