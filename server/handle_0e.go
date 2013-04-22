package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0E handles incoming requests of packet 0x0E: PlayerAction
func Handle0E(s *Server, player *player.Player) {
	p := new(packet.PlayerAction)
	p.ReadFrom(player.Conn)

	log.Printf("PlayerAction: %+v", p)
}
