package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0A handles incoming requests of packet 0x0A: Player
func Handle0A(s *Server, player *player.Player) {
	p := new(packet.Player)
	p.ReadFrom(player.Conn)

	log.Printf("Player: %+v", p)
}
