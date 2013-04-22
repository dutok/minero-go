package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle82 handles incoming requests of packet 0x82: SignUpdate
func Handle82(s *Server, player *player.Player) {
	p := new(packet.SignUpdate)
	p.ReadFrom(player.Conn)

	log.Printf("SignUpdate: %+v", p)
}
