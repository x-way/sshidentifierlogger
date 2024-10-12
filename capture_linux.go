// +build linux

package main

import (
	"github.com/gopacket/gopacket/pcapgo"
)

func LiveHandle(iface string) (*pcapgo.EthernetHandle, error) {
	return pcapgo.NewEthernetHandle(iface)
}
