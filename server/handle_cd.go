package server

import (
	"log"
	"time"

	"github.com/toqueteos/minero/proto/packet"
	"github.com/toqueteos/minero/server/player"
)

// HandleCD handles incoming requests of packet 0xCD: ClientStatuses
func HandleCD(s *Server, player *player.Player) {
	p := new(packet.ClientStatuses)
	p.ReadFrom(player.Conn)

	log.Println("HandleCD Payload:", p.Payload)
	switch p.Payload {
	case 0:
		// Send Login packet
		r := packet.LoginInfo{
			Entity:     33,
			LevelType:  "default",
			GameMode:   1,
			Dimension:  0,
			Difficulty: 2,
			MaxPlayers: 32,
		}
		r.WriteTo(player.Conn)

		// Save player login timestamp
		player.Since = time.Now().Unix()
	case 1:
	default:
		r := packet.Disconnect{"Weird packet 0xCB payload"}
		r.WriteTo(player.Conn)
	}

	r := packet.Disconnect{"Not implemented."}
	r.WriteTo(player.Conn)
}
