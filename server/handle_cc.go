package server

import (
	"log"

	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// HandleCC handles incoming requests of packet 0xCC: ClientSettings
func HandleCC(server *Server, sender *player.Player) {
	pkt := new(packet.ClientSettings)
	pkt.ReadFrom(sender.Conn)

	log.Printf("ClientSettings: %+v", pkt)
}
