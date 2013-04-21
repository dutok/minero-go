package server

import (
	"fmt"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/proto/ping"
	"github.com/toqueteos/minero/server/player"
)

// HandleFE handles incoming requests of packet 0xFE: ServerListPing
func HandleFE(s *Server, player *player.Player) {
	r := new(packet.ServerListPing)
	r.ReadFrom(player.Conn)

	if r.Magic != 1 {
		s := "Invalid %#x packet. Field Magic should be 1, got %d."
		reason := fmt.Sprintf(s, r.Id(), r.Magic)
		resp := packet.Disconnect{reason}
		resp.WriteTo(player.Conn)
		return
	}

	in := fmt.Sprintf("%d", s.in)
	max := fmt.Sprintf("%d", s.max)
	resp := ping.Ping(ping.Prepare(s.Motd(), in, max))
	resp.WriteTo(player.Conn)
}
