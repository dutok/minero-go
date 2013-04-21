package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle12 handles incoming requests of packet 0x12: Animation
func Handle12(s *Server, player *player.Player) {
	p := new(packet.Animation)
	p.ReadFrom(player.Conn)

	log.Printf("Animation: %+v", p)
}
