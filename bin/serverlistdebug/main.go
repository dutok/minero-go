// Fake Minecraft server list server. Doesn't require any other packages, this
// may change in the future.
package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("sld> ")
	log.SetFlags(log.Ltime)

	flag.Parse()
}

func main() {
	switch flag.NArg() {
	case 2:
		switch flag.Arg(0) {
		case "client":
			Client(flag.Arg(1))
		case "server":
			Server(flag.Arg(1))
		}
	default:
		log.Fatalln("Usage: serverlistdebug [client|server] [addr|port]")
	}
}
