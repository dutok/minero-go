package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0C handles incoming requests of packet 0x0C: PlayerLook
func Handle0C(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerLook)
	pkt.ReadFrom(sender.Conn)

	resp := &packet.EntityLook{
		Entity: sender.Id(),
		Yaw:    pkt.Yaw,
		Pitch:  pkt.Pitch,
	}
	server.BroadcastOthers(sender, resp)
}
