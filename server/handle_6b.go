package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle6B handles incoming requests of packet 0x6B: CreativeInventoryAction
func Handle6B(s *Server, player *player.Player) {
	p := new(packet.CreativeInventoryAction)
	p.ReadFrom(player.Conn)

	if p.Item != nil {
		log.Printf("CreativeInventoryAction: %+v", p)
	} else {
		log.Println("Slot is nil")
	}
}
