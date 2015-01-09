package server

import (
	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
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
	server.BroadcastPacket(resp)
}
