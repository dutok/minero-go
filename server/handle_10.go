package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle10 handles incoming requests of packet 0x10: ItemHeldChange
func Handle10(server *Server, sender *player.Player) {
	pkt := new(packet.ItemHeldChange)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ItemHeldChange: %+v", pkt)
}
