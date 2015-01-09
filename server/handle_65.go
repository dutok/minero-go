package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle65 handles incoming requests of packet 0x65: WindowClose
func Handle65(server *Server, sender *player.Player) {
	pkt := new(packet.WindowClose)
	pkt.ReadFrom(sender.Conn)

	log.Printf("WindowClose: %+v", pkt)
}
