package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle6A handles incoming requests of packet 0x6A: ConfirmTransaction
func Handle6A(server *Server, sender *player.Player) {
	pkt := new(packet.ConfirmTransaction)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ConfirmTransaction: %+v", pkt)
}
