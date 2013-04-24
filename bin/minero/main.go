package main

import (
	"log"

	"github.com/toqueteos/minero/server"
)

func main() {
	log.SetPrefix("minero> ")
	log.SetFlags(log.Ltime)

	s := server.New(nil)
	s.Run()
}
