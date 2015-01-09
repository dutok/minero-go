package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle07 handles incoming requests of packet 0x07: EntityInteract
func Handle07(server *Server, sender *player.Player) {
	pkt := new(packet.EntityInteract)
	pkt.ReadFrom(sender.Conn)

	log.Printf("EntityInteract: %+v", pkt)
}
