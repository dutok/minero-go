package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCB handles incoming requests of packet 0xCB: TabComplete
func HandleCB(s *Server, player *player.Player) {
	p := new(packet.TabComplete)
	p.ReadFrom(player.Conn)

	log.Printf("TabComplete: %+v", p)
}
