// +build !windows

package util

import (
  "syscall"
  "unsafe"

	"github.com/jcelliott/lumber"
)

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func GetTerminalSize() (int, int) {
  ws := &winsize{}
  retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
    uintptr(syscall.Stdin),
    uintptr(syscall.TIOCGWINSZ),
    uintptr(unsafe.Pointer(ws)))

  if int(retCode) == -1 {
  	lumber.Error("GetTerminalSize(): %s", errno.Error())
    return 30, 80
  }

  return int(ws.Row), int(ws.Col)
}