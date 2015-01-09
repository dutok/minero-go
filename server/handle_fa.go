package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// HandleFA handles incoming requests of packet 0xFA: PluginMessage
func HandleFA(server *Server, sender *player.Player) {
	pkt := new(packet.PluginMessage)
	pkt.ReadFrom(sender.Conn)

	log.Printf("PluginMessage: %+v", pkt)
}
