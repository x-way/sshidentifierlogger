package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
)

type logEntry struct {
	SrcIP    string `json:"src"`
	SrcPort  int    `json:"sport"`
	DstIP    string `json:"dst"`
	DstPort  int    `json:"dport"`
	Identity string `json:"sshid"`
	Time     string `json:"@timestamp"`
	Version  int    `json:"@version"`
}

var (
	iface  = flag.String("i", "eth0", "Interface to read packets from (only supported on Linux)")
	fname  = flag.String("r", "", "Filename to read from, overrides -i")
	server = flag.String("s", "", "Server IP, packets coming from this address are ignored (use commas to specify multiple IPs)")
	debug  = flag.Bool("d", false, "debug mode, log to stdout")
	logger syslog.Writer
)

func main() {
	logger, _ := syslog.New(syslog.LOG_DAEMON|syslog.LOG_INFO, "sshversionlogger")
	flag.Parse()

	if *fname != "" {
		if *debug {
			fmt.Println("reading packets from file " + *fname)
		}
		f, _ := os.Open(*fname)
		defer f.Close()
		if handle, err := pcapgo.NewReader(f); err != nil {
			if logerr := logger.Err(fmt.Sprintf("NewReader error: %v", err)); logerr != nil {
				log.Fatal(logerr)
			}
			log.Fatal("NewReader error:", err)
		} else {
			run(handle, handle.LinkType())
		}
	} else {
		if *debug {
			fmt.Println("reading packets from interface " + *iface)
		}
		if handle, err := liveHandle(*iface); err != nil {
			if logerr := logger.Err(fmt.Sprintf("NewEthernetHandle error: %v", err)); logerr != nil {
				log.Fatal(logerr)
			}
			log.Fatal("NewEthernetHandle error:", err)
		} else {
			defer handle.Close()
			run(handle, layers.LayerTypeEthernet)
		}
	}
}

func run(src gopacket.PacketDataSource, dec gopacket.Decoder) {
	source := gopacket.NewPacketSource(src, dec)
	source.Lazy = true
	source.NoCopy = true
	source.DecodeStreamsAsDatagrams = true
loop:
	for packet := range source.Packets() {
		if packet.Metadata().CaptureInfo.Length < 67 {
			continue
		}
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		if int(tcp.SrcPort) != 22 && int(tcp.DstPort) != 22 {
			continue
		}
		if !bytes.HasPrefix(tcp.BaseLayer.Payload, []byte("SSH-")) {
			continue
		}
		src, dst := packet.NetworkLayer().NetworkFlow().Endpoints()
		srcString := src.String()
		dstString := dst.String()
		ips := strings.Split(*server, ",")
		for i := range ips {
			if len(strings.TrimSpace(ips[i])) > 0 {
				if ips[i] == srcString {
					continue loop
				}
			}
		}

		lines := strings.Split(string(tcp.BaseLayer.Payload), "\n")
		entry := logEntry{
			SrcIP:    srcString,
			DstIP:    dstString,
			SrcPort:  int(tcp.SrcPort),
			DstPort:  int(tcp.DstPort),
			Identity: strings.TrimSuffix(lines[0], "\r"),
			Time:     packet.Metadata().Timestamp.UTC().Format(time.RFC3339),
			Version:  1,
		}

		if b, err := json.Marshal(entry); err == nil {
			if *debug {
				fmt.Println(string(b))
			} else {
				if err := logger.Info(string(b)); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
