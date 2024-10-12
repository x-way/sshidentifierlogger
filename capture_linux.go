//go:build linux
// +build linux

package main

import (
	"github.com/gopacket/gopacket/pcapgo"
)

func liveHandle(iface string) (*pcapgo.EthernetHandle, error) {
	return pcapgo.NewEthernetHandle(iface)
}
