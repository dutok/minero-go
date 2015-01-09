// Minecraft Server implementation in Go.
package main

import (
	"log"

	"github.com/minero/minero/server"
)

func main() {
	log.SetPrefix("minero> ")
	log.SetFlags(log.Ltime)

	s := server.New(nil)
	s.Run()
}
