package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var (
	PrintVersion = flag.Bool("v", false, "Display build info and exit")
	Address      = flag.String("a", "238.6.0.2:6666", "Address to cast to")

	Version   = "development"
	GitCommit = "development"
	BuildTime = "development"
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

	conn, err := NewUDPConn(*Address)
	if err != nil {
		log.Fatal(":( ", err)
	}
	defer conn.Close()
	p := ipv4.NewPacketConn(conn)
	p.SetMulticastTTL(3)

	for {
		t := fmt.Sprintf(time.Now().Format(time.RFC3339))
		log.Printf("sending: %s", t)
		conn.Write([]byte(t))
		time.Sleep(1 * time.Second)
	}
}
