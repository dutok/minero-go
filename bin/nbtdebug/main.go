package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/toqueteos/minero/proto/nbt"
)

const (
	// Uncompressed
	Servers = `C:\Users\Carlos\AppData\Roaming\.minecraft\servers.dat`
	// Gzip
	BigTest = "../_testdata/bigtest.nbt"
	Test    = "../_testdata/test.nbt"
)

func usage() {
	fmt.Println("Usage: nbtdebug nbtfile")
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return
	}

	var filename = flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Couldn't open file: %q.\n", filename)
	}

	c, err := nbt.Read(file)
	if err != nil {
		log.Fatalln("nbt.Read:", err)
	}

	for k, v := range c.Tags {
		parts := strings.SplitN(k, " ", 2)
		fmt.Printf("%q: %v\n", parts[0], v)
	}
}
