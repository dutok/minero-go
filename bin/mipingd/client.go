package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"

	"github.com/minero/minero/proto/packet"
)

func Client(addr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	log.Println("Connected to:", c.RemoteAddr())

	// Send request
	p := packet.ServerListPing{Magic: 1}
	p.WriteTo(c)

	// Read server response
	var buf bytes.Buffer
	io.Copy(&buf, c)

	// Read packet id
	id, _ := buf.ReadByte()

	if id != packet.PacketDisconnect {
		log.Fatalln("Unexpected packet id:", id)
	}

	r := new(packet.Disconnect)
	r.ReadFrom(&buf)

	s := strings.Split(r.Reason, "\x00")

	log.Println("Protocol version:", s[1])
	log.Println("Server version:", s[2])
	log.Println("MOTD:", s[3])
	log.Printf("Players: %s/%s\n", s[4], s[5])
}
