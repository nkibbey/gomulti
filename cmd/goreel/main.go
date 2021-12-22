package main

import (
	"encoding/hex"
	"log"
	"net"
)

const (
	maxDSz = 8192
)

func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

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
