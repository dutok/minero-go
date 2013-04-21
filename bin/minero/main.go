package main

import (
	"log"

	"github.com/toqueteos/minero/server"
)

func main() {
	log.SetPrefix("minero> ")
	log.SetFlags(log.Ltime)

	s := server.New("", "25600")
	s.Run()
}
