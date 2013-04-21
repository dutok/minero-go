package server

import (
	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle0D handles incoming requests of packet 0x0D: PlayerPosLook
func Handle0D(s *Server, player *player.Player) {
	p := new(packet.PlayerPosLook)
	p.ReadFrom(player.Conn)

	if p.Y > p.Stance {
		r := &packet.Disconnect{"Weird packet 0x0D. Server didn't switch Y with Stance."}
		r.WriteTo(player.Conn)
	}
}
