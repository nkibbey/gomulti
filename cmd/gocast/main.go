package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"

)

func NewUDPConn(addr string) (*net.UDPConn, error) {
	a, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, a)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	conn, err := NewUDPConn("238.6.0.2:6666")
	if err != nil {
		log.Fatal(":( ", err)
	}
	defer conn.Close()
	p := ipv4.NewPacketConn(conn)
	p.SetMulticastTTL(3)

	for {
		t := fmt.Sprintf("%v,\n", time.Now())
		log.Printf("sending: %s", t)
		conn.Write([]byte(t))
		time.Sleep(1 * time.Second)
	}
}
