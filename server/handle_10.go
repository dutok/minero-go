package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle10 handles incoming requests of packet 0x10: ItemHeldChange
func Handle10(server *Server, sender *player.Player) {
	pkt := new(packet.ItemHeldChange)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ItemHeldChange: %+v", pkt)
}
