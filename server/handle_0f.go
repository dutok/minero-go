package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0F handles incoming requests of packet 0x0F: PlayerBlockPlace
func Handle0F(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerBlockPlace)
	pkt.ReadFrom(sender.Conn)

	log.Printf("PlayerBlockPlace: %+v", pkt)
}
