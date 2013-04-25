// Server proxy and play session log recorder.
package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var Flags struct {
	Src, Dst, File string
	Record, Gzip   bool
}

// Usage of each flag
var u = map[string]string{
	"src":    "Proxy destination address.",
	"dst":    "Proxy listen address.",
	"file":   "Stream storage file.",
	"record": "Save stream to storage?",
	"gzip":   "Compress stream with gzip before storing?",
}

func init() {
	flag.StringVar(&Flags.Src, "src", "127.0.0.1:25565", u["src"])
	flag.StringVar(&Flags.Dst, "dst", "127.0.0.1:26665", u["dst"])
	flag.StringVar(&Flags.File, "file", "proxy.log", u["file"])
	flag.BoolVar(&Flags.Record, "record", true, u["record"])
	flag.BoolVar(&Flags.Gzip, "gzip", true, u["gzip"])

	flag.Parse()
}

func main() {
	log.SetPrefix("miproxy> ")
	log.SetFlags(log.Ltime)

	var err error

	log.Println("Pinging remote server.")
	serverConn, err := net.Dial("tcp", Flags.Src)
	if err != nil {
		log.Println("Remote server: OFFLINE.")
		log.Fatalln(err)
	}
	defer serverConn.Close()
	log.Println("Remote server: ONLINE.")

	log.Println("Setting up proxy.")
	listener, err := net.Listen("tcp", Flags.Dst)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Println("Done!")
	log.Printf("Listening on %q\n", Flags.Dst)

	var file io.Writer

	// Log proxy stream
	if Flags.Record {
		var err error

		fileFlags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		file, err = os.OpenFile(Flags.File, fileFlags, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		if Flags.Gzip {
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
