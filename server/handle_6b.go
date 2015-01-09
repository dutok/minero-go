package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
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
