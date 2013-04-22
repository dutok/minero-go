package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle00 handles incoming requests of packet 0x00: KeepAlive
func Handle00(s *Server, player *player.Player) {
	p := new(packet.KeepAlive)
	p.ReadFrom(player.Conn)

	log.Printf("KeepAlive: %+v", p)
}
