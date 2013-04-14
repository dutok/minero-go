// Server list ping client & server.
package main

import (
	"flag"
	"log"
)

var Flags [6]string

func init() {
	log.SetPrefix("mipingd> ")
	log.SetFlags(log.Ltime)

	Flags[0] = "§1"
	flag.StringVar(&Flags[1], "proto", "60", "")
	flag.StringVar(&Flags[2], "server", "1.5", "")
	flag.StringVar(&Flags[3], "motd", "§9minero§r Server", "")
	flag.StringVar(&Flags[4], "n", "0", "")
	flag.StringVar(&Flags[5], "m", "64", "")
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
