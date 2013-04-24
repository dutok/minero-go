package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0D handles incoming requests of packet 0x0D: PlayerPosLook
func Handle0D(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerPosLook)
	pkt.ReadFrom(sender.Conn)

	if pkt.Y > pkt.Stance {
		r := &packet.Disconnect{"Weird packet 0x0D. Server didn't switch Y with Stance."}
		r.WriteTo(sender.Conn)
		return
	}

	// Player is ready for broadcasts
	sender.SetReady()
	// Save position
	sender.SetPos(pkt.X, pkt.Y, pkt.Z)
	sender.SetLook(pkt.Pitch, pkt.Yaw)

	var p packet.Packet

	// BUG(toqueteos): NMS handle relative movements with other packet
	// `packet.EntityRelMove` right now we just teleport to the destination.
	p = &packet.EntityTeleport{
		Entity: sender.Id(),
		X:      pkt.X,
		Y:      pkt.Y,
		Z:      pkt.Z,
		Yaw:    pkt.Yaw,
		Pitch:  pkt.Pitch,
	}
	server.BroadcastOthers(sender, p)

	p = &packet.EntityHeadLook{
		Entity:  sender.Id(),
		HeadYaw: pkt.Yaw,
	}
	server.BroadcastOthers(sender, p)
}
