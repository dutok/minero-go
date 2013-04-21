package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle13 handles incoming requests of packet 0x13: EntityAction
func Handle13(s *Server, player *player.Player) {
	p := new(packet.EntityAction)
	p.ReadFrom(player.Conn)

	log.Printf("EntityAction: %+v", p)
}
