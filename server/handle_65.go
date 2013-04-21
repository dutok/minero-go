package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle65 handles incoming requests of packet 0x65: WindowClose
func Handle65(s *Server, player *player.Player) {
	p := new(packet.WindowClose)
	p.ReadFrom(player.Conn)

	log.Printf("WindowClose: %+v", p)
}
