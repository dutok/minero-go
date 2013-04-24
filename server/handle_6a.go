package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle6A handles incoming requests of packet 0x6A: ConfirmTransaction
func Handle6A(server *Server, sender *player.Player) {
	pkt := new(packet.ConfirmTransaction)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ConfirmTransaction: %+v", pkt)
}
