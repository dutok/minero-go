package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Flags struct {
	Src, Dst, File string
	Record, Gzip   bool
}

// Value of each flag
var f Flags

// Usage of each flag
var u = map[string]string{
	"src":    "Address where proxy will read from.",
	"dst":    "Address where proxy will listen for.",
	"file":   "File to store proxy stream.",
	"record": "Save proxy stream to storage?",
	"gzip":   "Compress with gzip storage?",
}

func init() {
	flag.StringVar(&f.Src, "src", "127.0.0.1:25565", u["src"])
	flag.StringVar(&f.Dst, "dst", "127.0.0.1:26665", u["dst"])
	flag.StringVar(&f.File, "file", "proxy.log", u["file"])
	flag.BoolVar(&f.Record, "record", true, u["record"])
	flag.BoolVar(&f.Gzip, "gzip", true, u["gzip"])

	flag.Parse()
}

func main() {
	var err error

	fmt.Print("Is remote server online? ... ")
	serverConn, err := net.Dial("tcp", f.Src)
	if err != nil {
		log.Fatalln(err)
	}
	defer serverConn.Close()
	fmt.Println("OK")

	fmt.Print("Setting up proxy ... ")
	listener, err := net.Listen("tcp", f.Dst)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	fmt.Println("OK")
	fmt.Printf("Listening on %q. OK\n", f.Dst)

	var file io.Writer

	// Log proxy stream
	if f.Record {
		var err error

		fileFlags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		file, err = os.OpenFile(f.File, fileFlags, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		if f.Gzip {
			file = gzip.NewWriter(file)
		}
	}

	var client, server io.Reader
	if file != nil {
		// Log whatever we read from remote server
		server = io.TeeReader(serverConn, &PrefixWriter{file, "S"})
	} else {
		server = serverConn
	}

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		// Log whatever we read from remote server
		client = io.TeeReader(clientConn, &PrefixWriter{file, "C"})

		// Server <- Client
		go io.Copy(serverConn, client)

		// Client <- Server
		go io.Copy(clientConn, server)
	}
}

type PrefixWriter struct {
	rw     io.Writer
	Prefix string
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	var nn int
	prefix := []byte(pw.Prefix + " ")

	// Write prefix
	n, err = pw.rw.Write(prefix)
	if err != nil {
		return
	}
	nn += n

	// Write content
	n, err = pw.rw.Write(p)
	if err != nil {
		// Prefix write succesful
		n = nn
		return
	}
	nn += n

	return nn, nil
}
