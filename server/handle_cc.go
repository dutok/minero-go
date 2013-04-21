package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCC handles incoming requests of packet 0xCC: ClientSettings
func HandleCC(s *Server, player *player.Player) {
	p := new(packet.ClientSettings)
	p.ReadFrom(player.Conn)

	log.Printf("ClientSettings: %+v", p)
}
