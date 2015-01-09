package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle13 handles incoming requests of packet 0x13: EntityAction
func Handle13(server *Server, sender *player.Player) {
	pkt := new(packet.EntityAction)
	pkt.ReadFrom(sender.Conn)

	log.Printf("EntityAction: %+v", pkt)
}
