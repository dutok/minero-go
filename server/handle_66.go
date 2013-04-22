package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle66 handles incoming requests of packet 0x66: WindowClick
func Handle66(s *Server, player *player.Player) {
	p := new(packet.WindowClick)
	p.ReadFrom(player.Conn)

	log.Printf("WindowClick: %+v", p)
}
