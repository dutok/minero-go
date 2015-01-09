package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle66 handles incoming requests of packet 0x66: WindowClick
func Handle66(server *Server, sender *player.Player) {
	pkt := new(packet.WindowClick)
	pkt.ReadFrom(sender.Conn)

	log.Printf("WindowClick: %+v", pkt)
}
