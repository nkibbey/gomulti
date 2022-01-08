package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var (
	PrintVersion = flag.Bool("v", false, "Display build info and exit")
	Address      = flag.String("a", "238.6.0.2:6666", "Address to cast to")

	Version   = "development"
	GitCommit = "development"
	BuildTime = "development"
)

const (
	maxDSz = 8192
)

func joinBroadcastInterfaces(p *ipv4.PacketConn, addr *net.UDPAddr) {
	intfs, err := net.Interfaces()
	if err != nil {
		log.Println("couldn't get interfaces for joining,", err)
	}

	for _, i := range intfs {
		if i.Name == "lo" { // doesn't seem to show up anyways
			continue
		} else if i.Flags&net.FlagMulticast == 0 {
			log.Println("skipping", i.Name, "because intf does not have multicast ->", i.Flags)
			continue
		}

		p.LeaveGroup(&i, addr) // in case it was already joined in net.listenmulticastudp
		log.Println("joining", i.Name, ", index", i.Index, "...")
		if err := p.JoinGroup(&i, addr); err != nil {
			log.Println("group join err", err)
		}
	}
}

func Listen(address string) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	joinBroadcastInterfaces(p, addr)

	conn.SetReadBuffer(maxDSz)

	if err := p.SetControlMessage(ipv4.FlagInterface, true); err != nil {
		log.Println("couldn't set cm for receive interface info because", err)
	}

	readConn(p, maxDSz)
	p.Close()
}

func readConn(p *ipv4.PacketConn, bufSz int) {
	for {
		buf := make([]byte, bufSz)
		n, cm, src, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal("readfrom fail", err)
		}

		if cm != nil {
			log.Println(n, "bytes read from", src, "on interface:", cm.IfIndex)
		} else { // shouldn't be hit if setcontrolmessage worked
			log.Println(n, "bytes read from", src)
		}

		if len(buf) > 8 {
			log.Printf("%s", buf[8:])
		} else {
			log.Println("UNKNOWN UDP DATAGRAM: ", hex.Dump(buf[:n]))
		}
	}
}

// buildInfo provides information about this build
func buildInfo() string {
	return fmt.Sprintf("Version: %s\nBuild Time: %s\nGitCommit: %s", Version, BuildTime, GitCommit)
}

func main() {
	flag.Parse()
	if *PrintVersion {
		log.Printf("------BUILD INFO-----\n%s\n-----------------------------------------", buildInfo())
		return
	}

	Listen(*Address)
}
