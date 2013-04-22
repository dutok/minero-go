package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle07 handles incoming requests of packet 0x07: EntityInteract
func Handle07(s *Server, player *player.Player) {
	p := new(packet.EntityInteract)
	p.ReadFrom(player.Conn)

	log.Printf("EntityInteract: %+v", p)
}
