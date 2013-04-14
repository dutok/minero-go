// NBT pretty printer.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing/iotest"

	"github.com/toqueteos/minero/proto/nbt"
)

var Flags struct {
	Quiet bool
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("minbtd> ")
	flag.BoolVar(&Flags.Quiet, "q", false, "")
	flag.Parse()
}

func main() {
	if flag.NArg() != 1 {
		fmt.Println("Usage: minbtd file")
		return
	}

	var (
		filename = flag.Arg(0)
		err      error
	)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Couldn't open file: %q.\n", filename)
	}

	var (
		r, rr io.Reader
		gzip  bool
	)

	if !Flags.Quiet {
		r = iotest.NewReadLogger("r:", file)
	} else {
		r = file
	}

	rr, gzip, err = nbt.GuessCompression(r)
	if err != nil {
		log.Fatalf("Error guessing file compression: %q\n", err)
	}

	if gzip {
		log.Println("Detected gzip file.")
	}

	c, err := nbt.ReadRaw(rr)
	if err != nil {
		log.Fatalln("nbt.Read:", err)
	}

	log.Println("Top level compound name:", c.Name)
	for k, v := range c.Value {
		parts := strings.SplitN(k, " ", 2)
		log.Printf("%q: %v\n", parts[0], v)
	}
}
