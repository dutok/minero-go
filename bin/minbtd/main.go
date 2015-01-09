// NBT pretty printer.
package main

import (
	"compress/gzip"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing/iotest"

	"github.com/minero/minero/proto/nbt"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("minbtd> ")

	flag.BoolVar(&Flags.Quiet, "q", false, "")
	flag.Parse()
}

const (
	None = iota
	Gzip
	Zlib
)

var Flags struct {
	Quiet bool
}

func main() {
	switch flag.NArg() {
	case 1:
		Prepare(flag.Arg(0), "gzip")
	case 2:
		Prepare(flag.Arg(0), flag.Arg(1))
	default:
		fmt.Println("Usage: minbtd file [compression]")
		fmt.Println("Default compression: GZIP")
	}
}

func Prepare(file string, mode string) {
	var f io.ReadCloser
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Couldn't open file: %q.\n", flag.Arg(0))
	}
	defer f.Close()

	switch mode {
	case "none":
		Debug(f)
	case "gzip":
		r, err := gzip.NewReader(f)
		if err != nil {
			log.Fatalln(err)
		}
		defer r.Close()
		Debug(r)
	case "zlib":
		r, err := zlib.NewReader(f)
		if err != nil {
			log.Fatalln(err)
		}
		defer r.Close()
		Debug(r)
	default:
		log.Fatalf("Unknown compression scheme %q.\n", mode)
	}
}

func Debug(f io.Reader) {
	var r io.Reader
	if !Flags.Quiet {
		r = iotest.NewReadLogger("r:", f)
	} else {
		r = f
	}

	c, err := nbt.Read(r)
	if err != nil {
		log.Fatalln("nbt.Read:", err)
	}

	log.Printf("Top level compound name: %q\n", c.Name)
	fmt.Println(nbt.Pretty(c.String()))
}
