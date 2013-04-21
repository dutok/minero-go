package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0B handles incoming requests of packet 0x0B: PlayerPos
func Handle0B(s *Server, player *player.Player) {
	p := new(packet.PlayerPos)
	p.ReadFrom(player.Conn)

	// log.Printf("PlayerPos: %+v", p)
}
