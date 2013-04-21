package server

import (
	"log"
	"strings"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// Handle03 handles incoming requests of packet 0x03: ChatMessage
func Handle03(s *Server, player *player.Player) {
	p := new(packet.ChatMessage)
	p.ReadFrom(player.Conn)

	log.Printf("ChatMessage: %+v", p)

	if strings.HasPrefix(p.Message, "/") {
		HandleCommand(player, p.Message[1:])
	} else {
		// ... send to all other players
	}
}

func HandleCommand(player *player.Player, m string) {
	f := strings.Fields(m)

	switch f[0] {
	case "gamemode":
		r := packet.GameStateChange{Reason: 3}

		switch f[1] {
		case "0", "s":
			r.GameMode = 0
		case "1", "c":
			r.GameMode = 1
		case "2", "a":
			r.GameMode = 2
		}

		r.WriteTo(player.Conn)
	default:
		log.Println("")
	}
}
