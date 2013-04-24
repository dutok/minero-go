package server

import (
	"fmt"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/proto/ping"
	"github.com/toqueteos/minero/server/player"
)

// HandleFE handles incoming requests of packet 0xFE: ServerListPing
func HandleFE(server *Server, sender *player.Player) {
	pkt := new(packet.ServerListPing)
	pkt.ReadFrom(sender.Conn)

	if pkt.Magic != 1 {
		s := "Invalid %#x packet. Field Magic should be 1, got %d."
		reason := fmt.Sprintf(s, pkt.Id(), pkt.Magic)
		resp := packet.Disconnect{reason}
		resp.WriteTo(sender.Conn)
		return
	}

	in := fmt.Sprintf("%d", len(server.playerList))
	max := server.config.Get("server.max_players")
	resp := ping.Ping(ping.Prepare(server.Motd, in, max))
	resp.WriteTo(sender.Conn)
}
