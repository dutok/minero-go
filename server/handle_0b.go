package server

import (
	"github.com/minero/minero/proto/packet"
	"github.com/minero/minero/server/player"
)

// Handle0B handles incoming requests of packet 0x0B: PlayerPos
func Handle0B(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerPos)
	pkt.ReadFrom(sender.Conn)
	// server.BroadcastPacket(pkt)
}
