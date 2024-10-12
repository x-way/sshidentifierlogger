// +build !linux

package main

import (
	"errors"

	"github.com/gopacket/gopacket"
)

var errNotLinux = errors.New("not implemented for OS other than linux")

type FakeHandle struct {

	// Dummy
}

func (f FakeHandle) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	// Dummy
	return nil, gopacket.CaptureInfo{}, errNotLinux
}

func (f *FakeHandle) Close() {
	// Dummy
}

func LiveHandle(_ string) (FakeHandle, error) {
	return FakeHandle{}, errNotLinux
}
