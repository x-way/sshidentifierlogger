//go:build !linux
// +build !linux

package main

import (
	"errors"

	"github.com/gopacket/gopacket"
)

var errNotLinux = errors.New("not implemented for OS other than linux")

type fakeHandle struct {
	// Dummy
}

func (f fakeHandle) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	// Dummy
	return nil, gopacket.CaptureInfo{}, errNotLinux
}

func (f *fakeHandle) Close() {
	// Dummy
}

func liveHandle(_ string) (fakeHandle, error) {
	return fakeHandle{}, errNotLinux
}
