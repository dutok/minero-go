package main

import (
	"log"
	"net"
	"strings"

	packet "github.com/toqueteos/minero/proto/minecraft"
	"github.com/toqueteos/minero/proto/minecraft/ping"
)

func Server(addr string) {
	log.SetPrefix("serverlistdebug> ")
	log.SetFlags(log.Ltime)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		// Serve multiple connections concurrently.
		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer c.Close()

	log.Println("Got connection from:", c.RemoteAddr())

	// Read first two bytes NMC sends
	var buf = make([]byte, 2)
	_, err := c.Read(buf)
	if err != nil {
		log.Printf("%s: %s\n", c.RemoteAddr(), err)
		return
	}

	// Equal to 0xFE?
	if buf[0] != packet.PacketServerListPing || buf[1] != 0x01 {
		return
	}

	// Send response
	p := packet.Disconnect{
		Reason: strings.Join(Flags[:], "\x00\x00"),
	}
	_ = p
	// p.WriteTo(c)

	c.Write(ping.Ping(Flags[:]))
}
