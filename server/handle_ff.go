package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleFF handles incoming requests of packet 0xFF: Disconnect
func HandleFF(s *Server, player *player.Player) {
	p := new(packet.Disconnect)
	p.ReadFrom(player.Conn)

	log.Printf("Player %q exit. Reason: %s", player.Name, p.Reason)
}
