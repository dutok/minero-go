package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCA handles incoming requests of packet 0xCA: PlayerAbilities
func HandleCA(s *Server, player *player.Player) {
	p := new(packet.PlayerAbilities)
	p.ReadFrom(player.Conn)

	log.Printf("PlayerAbilities: %+v", p)
}
