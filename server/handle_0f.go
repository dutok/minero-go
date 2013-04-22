package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0F handles incoming requests of packet 0x0F: PlayerBlockPlace
func Handle0F(s *Server, player *player.Player) {
	p := new(packet.PlayerBlockPlace)
	p.ReadFrom(player.Conn)

	log.Printf("PlayerBlockPlace: %+v", p)
}
