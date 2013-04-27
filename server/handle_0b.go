package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0B handles incoming requests of packet 0x0B: PlayerPos
func Handle0B(server *Server, sender *player.Player) {
	pkt := new(packet.PlayerPos)
	pkt.ReadFrom(sender.Conn)
	// server.BroadcastPacket(pkt)
}
