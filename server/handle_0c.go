package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0C handles incoming requests of packet 0x0C: PlayerLook
func Handle0C(s *Server, player *player.Player) {
	p := new(packet.PlayerLook)
	p.ReadFrom(player.Conn)

	log.Printf("PlayerLook: %+v", p)
}
