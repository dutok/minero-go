package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle10 handles incoming requests of packet 0x10: ItemHeldChange
func Handle10(s *Server, player *player.Player) {
	p := new(packet.ItemHeldChange)
	p.ReadFrom(player.Conn)

	log.Printf("ItemHeldChange: %+v", p)
}
