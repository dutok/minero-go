package server

import (
	"log"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleFA handles incoming requests of packet 0xFA: PluginMessage
func HandleFA(s *Server, player *player.Player) {
	p := new(packet.PluginMessage)
	p.ReadFrom(player.Conn)

	log.Printf("PluginMessage: %+v", p)
}
