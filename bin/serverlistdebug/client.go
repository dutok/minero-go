package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"testing/iotest"
)

func Client(addr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to:", c.RemoteAddr())

	list := []byte{0xfe, 0x01}
	buf := bytes.NewBuffer(list)
	io.Copy(c, iotest.NewReadLogger("c->s", buf))

	// Read whatever server sends
	buf.Reset()
	io.Copy(iotest.NewWriteLogger("s->c", buf), c)

	c.Close()
}
