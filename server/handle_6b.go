package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle6B handles incoming requests of packet 0x6B: CreativeInventoryAction
func Handle6B(server *Server, sender *player.Player) {
	pkt := new(packet.CreativeInventoryAction)
	pkt.ReadFrom(sender.Conn)

	if pkt.Item != nil {
		log.Printf("CreativeInventoryAction: %+v", pkt)
	} else {
		log.Println("nil Slot")
	}
}
