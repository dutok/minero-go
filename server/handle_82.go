package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle82 handles incoming requests of packet 0x82: SignUpdate
func Handle82(server *Server, sender *player.Player) {
	pkt := new(packet.SignUpdate)
	pkt.ReadFrom(sender.Conn)

	log.Printf("SignUpdate: %+v", pkt)
}
