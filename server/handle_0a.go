package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0A handles incoming requests of packet 0x0A: Player
func Handle0A(server *Server, sender *player.Player) {
	pkt := new(packet.Player)
	pkt.ReadFrom(sender.Conn)

	resp := &packet.Entity{sender.Id()}
	server.BroadcastOthers(sender, resp)
}
