package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle6A handles incoming requests of packet 0x6A: ConfirmTransaction
func Handle6A(s *Server, player *player.Player) {
	p := new(packet.ConfirmTransaction)
	p.ReadFrom(player.Conn)

	log.Printf("ConfirmTransaction: %+v", p)
}
