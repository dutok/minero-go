package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle6C handles incoming requests of packet 0x6C: EnchantItem
func Handle6C(s *Server, player *player.Player) {
	p := new(packet.EnchantItem)
	p.ReadFrom(player.Conn)

	log.Printf("EnchantItem: %+v", p)
}
