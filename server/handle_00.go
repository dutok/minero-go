package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle00 handles incoming requests of packet 0x00: KeepAlive
func Handle00(server *Server, sender *player.Player) {
	pkt := new(packet.KeepAlive)
	pkt.ReadFrom(sender.Conn)

	log.Printf("KeepAlive: %+v", pkt)
}
