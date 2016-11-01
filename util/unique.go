package util

import (
	"crypto/md5"
	"fmt"
	"net"
)

func UniqueID() string {
	interfaces, _ := net.Interfaces()

	for _, i := range interfaces {
		if i.Flags&net.FlagLoopback != net.FlagLoopback {
			return fmt.Sprintf("%x", md5.Sum(i.HardwareAddr))
		}
	}

	return "generic"
}
