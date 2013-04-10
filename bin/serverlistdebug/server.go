package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"testing/iotest"
	"unicode/utf16"
)

var Ping = []string{
	"§1",  // Required
	"60",  // Protocol
	"1.5", // Server version
	"§9minero§r Server", // Motd
	"0",                 // Current players
	"32",                // Max players
}

func WriteBlob(ping []string) []byte {
	var (
		err  error
		sbuf bytes.Buffer
	)

	// Encode string so we can compute it's length
	sbuf.Write([]byte{0x00, 0xA7, 0x00, 0x31})

	for _, s := range ping[1:] {
		// NUL
		sbuf.WriteByte(0)
		sbuf.WriteByte(0)
		// Write UCS-2 string
		ucs2 := utf16.Encode([]rune(s))
		err = binary.Write(&sbuf, binary.BigEndian, ucs2)
		if err != nil {
			return nil
		}
	}

	// Write Disconnect packet Id (byte)
	buf := bytes.NewBuffer([]byte{0xff})
	// Write length of string (short)
	length := sbuf.Len() / 2
	err = binary.Write(buf, binary.BigEndian, uint16(length))
	if err != nil {
		return nil
	}
	buf.Write(sbuf.Bytes())

	return buf.Bytes()
}

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
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			log.Println("Got connection from:", c.RemoteAddr())

			// Read whatever client sends
			var rb bytes.Buffer
			io.CopyN(iotest.NewWriteLogger("s->c", &rb), c, 2)

			// Client sent [0xFE, 0x01]?
			if bytes.Equal(rb.Bytes(), []byte{0xFE, 0x01}) {
				buf := bytes.NewBuffer(WriteBlob(Ping))

				// Send response
				io.Copy(c, iotest.NewReadLogger("c->s", buf))
			}

			c.Close()
		}(conn)
	}
}
