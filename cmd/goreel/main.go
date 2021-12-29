package main

import (
	"encoding/hex"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

const (
	maxDSz = 8192
)

// listen2 attempt to join all interface groups
func listen2(conn *net.UDPConn, addr *net.UDPAddr) {
	p := ipv4.NewPacketConn(conn)
	intfs, err := net.Interfaces()
	if err != nil {
		log.Println("couldn't get interfaces for joining,", err)
	}

	for _, i := range intfs {
		if i.Name == "lo" { // doesn't seem to show up anyways
			continue
		}
		log.Println("joining", i.Name, "...")
		if err := p.JoinGroup(&i, addr); err != nil {
			log.Println("group join err", err)
		}
		// defer p.LeaveGroup(&i, &addr)
	}
}

func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	listen2(conn, addr)

	conn.SetReadBuffer(maxDSz)

	for {
		buf := make([]byte, maxDSz)
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("readudp fail", err)
		}

		handler(src, n, buf)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	if len(b) > 8 {
		log.Printf("%s", b[8:])
	} else {
		log.Println("UNKNOWN UDP DATAGRAM: ", hex.Dump(b[:n]))
	}
}

func main() {
	Listen("239.6.0.2:6666", msgHandler)
}
