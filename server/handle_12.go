package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle12 handles incoming requests of packet 0x12: Animation
func Handle12(server *Server, sender *player.Player) {
	pkt := new(packet.Animation)
	pkt.ReadFrom(sender.Conn)

	log.Printf("Animation: %+v", pkt)
}
